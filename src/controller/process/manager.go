package process

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/process"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type processData struct {
	process     process.Process
	cmd         *exec.Cmd
	isStart     bool
	isRunning   bool
	mutex       sync.Mutex
	reader      *bufio.Reader
	errReader   *bufio.Reader
	logFilePath string
	fileMutex   sync.Mutex
}

type manager struct {
	processList []*processData
	statistic   base.Statistic
	mutex       sync.Mutex
}

func stopProcess(pid int) error {
	pro, err := os.FindProcess(pid)
	if err != nil {
		logutils.Error("Failed to FindProcess. error: ", err)
		return err
	}
	err = pro.Kill()
	if err != nil {
		logutils.Error("Failed to Kill. error: ", err)
		return err
	}
	_, _ = pro.Wait()
	return nil
}

func (m *manager) run() {
	processList, err := process.GetProcessList()
	if err != nil {
		logutils.Error("Failed to GetProcessList. error: ", err)
		return
	}

	for _, p := range processList {
		data := new(processData)
		data.process = *p
		if data.process.Pid != 0 {
			_ = stopProcess(data.process.Pid)
			updateProcessStatus(data, false, 0)
		}
		if data.process.Enable {
			_ = m.start(data)
		} else {
			data.isStart = false
			updateProcessStatus(data, false, 0)
		}
		m.processList = append(m.processList, data)
	}
	_ = m.updateStatistic()
}

func updateProcessStatus(p *processData, started bool, pid int) {
	status := process.Status{
		ProcessId: p.process.Id,
		Type:      base.StatusTypeStarted,
	}
	if started {
		status.Value = "1"
		p.process.Pid = pid
		p.process.StartTime = time.Now().Local().Format("2006-01-02 15:04:05")
	} else {
		status.Value = "0"
		p.process.Pid = 0
	}
	_ = process.UpdateStatus(&status, nil)
	_ = process.UpdateProcess(p.process.Id, &p.process, nil)
}

func (m *manager) storeLog(p *processData, wg *sync.WaitGroup) (*os.File, error) {
	var file *os.File
	fileName := fmt.Sprintf("process-log/%s-%d.log", p.process.Name, p.cmd.Process.Pid)
	pwd, err := os.Getwd()
	if err != nil {
		logutils.Error("Failed to Getwd. error: ", err)
		return nil, err
	}
	filePath := filepath.Join(pwd, fileName)
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		logutils.Error("Failed to MkdirAll. error: ", err)
		return nil, err
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logutils.Error("Failed to OpenFile. error: ", err)
		return nil, err
	}
	p.logFilePath = filePath
	wg.Add(1)
	go func() {
		for p.reader != nil {
			line, err := p.reader.ReadString('\n')
			if err != nil {
				break
			}
			p.fileMutex.Lock()
			_, err = file.WriteString(line)
			p.fileMutex.Unlock()
			if err != nil {
				break
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for p.errReader != nil {
			line, err := p.errReader.ReadString('\n')
			if err != nil {
				break
			}
			p.fileMutex.Lock()
			_, err = file.WriteString(line)
			p.fileMutex.Unlock()
			if err != nil {
				break
			}
		}
		wg.Done()
	}()
	return file, nil
}

func (m *manager) start(p *processData) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isStart {
		return nil
	}
	p.isStart = true
	p.isRunning = true
	go func() {
		for p.isStart {
			time.Sleep(100 * time.Millisecond)
			p.cmd = new(exec.Cmd)
			p.cmd.Path = p.process.Path
			p.cmd.Args = []string{
				p.cmd.Path,
			}
			p.cmd.Args = append(p.cmd.Args, strings.Split(p.process.Config, " ")...)
			p.cmd.Dir = p.process.Dir

			for i, arg := range p.cmd.Args {
				logutils.TraceF("%s arg%d %s", p.process.Name, i+1, arg)
			}

			stdout, err := p.cmd.StdoutPipe()
			if err != nil {
				logutils.ErrorF("Failed to get %s StdoutPipe. error: %v", p.process.Name, err)
				p.cmd = nil
				continue
			} else {
				p.reader = bufio.NewReader(stdout)
			}

			stderr, err := p.cmd.StderrPipe()
			if err != nil {
				logutils.ErrorF("Failed ato get %s StderrPipe. error: %v", p.process.Name, err)
				p.cmd = nil
				continue
			} else {
				p.errReader = bufio.NewReader(stderr)
			}

			if err := p.cmd.Start(); err != nil {
				logutils.ErrorF("Failed to start %s. error: %v", p.process.Name, err)
				updateProcessStatus(p, false, 0)
				p.cmd = nil
				p.reader = nil
				continue
			}
			logutils.Trace("Success to start ", p.process.Name)
			updateProcessStatus(p, true, p.cmd.Process.Pid)

			var wg sync.WaitGroup
			file, _ := m.storeLog(p, &wg)

			if err := p.cmd.Wait(); err != nil {
				logutils.WarningF("%s quit. error: %v", p.process.Name, err)
			} else {
				logutils.WarningF("%s quit.", p.process.Name)
			}
			updateProcessStatus(p, false, 0)

			wg.Wait()
			if file != nil {
				_ = file.Close()
			}

			p.cmd = nil
			p.reader = nil
			p.errReader = nil
		}
		p.isRunning = false
	}()
	return nil
}

func (m *manager) stop(p *processData) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.isStart {
		return nil
	}
	p.isStart = false
	if err := p.cmd.Process.Kill(); err != nil {
		logutils.Error("Failed to Kill. error: ", err)
		return err
	}
	for i := 0; i < 3000; i++ {
		time.Sleep(10 * time.Millisecond)
		if p.isRunning == false {
			break
		}
	}
	return nil
}

func (m *manager) restart(p *processData) error {
	if p.isStart {
		_ = m.stop(p)
	}
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		status, err := process.GetStatus(p.process.Id, base.StatusTypeStarted)
		if err != nil {
			logutils.Error("Failed to GetStatus. error: ", err)
			break
		}
		if status.Value == "0" {
			break
		}
	}
	return m.start(p)
}

func (m *manager) getProcessData(id int) (*processData, error) {
	for _, p := range m.processList {
		if p.process.Id == id {
			return p, nil
		}
	}
	return nil, errors.New("进程信息不存在")
}

func (m *manager) updateStatistic() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.statistic.Total = 0
	m.statistic.Running = 0
	m.statistic.Stopped = 0
	m.statistic.Disable = 0
	for _, p := range m.processList {
		m.statistic.Total++
		if !p.process.Enable {
			m.statistic.Disable++
		} else {
			s, err := process.GetStatus(p.process.Id, base.StatusTypeStarted)
			if err != nil {
				m.statistic.Stopped++
			} else {
				if s.Value == "0" {
					m.statistic.Stopped++
				} else {
					m.statistic.Running++
				}
			}
		}
	}
	return nil
}

func (m *manager) quit() {
	for _, data := range m.processList {
		_ = m.stop(data)
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			status, err := process.GetStatus(data.process.Id, base.StatusTypeStarted)
			if err != nil {
				logutils.Error("Failed to GetStatus. error: ", err)
				break
			}
			if status.Value == "0" {
				break
			}
		}
	}
}

type processHandler struct {
	m *manager
}

func (h *processHandler) Handle(key int, value *bus.Resource) {
	isEnableChanged := false
	if key == base.KeyProcess {
		if value.Status == base.StatusUpdated {
			p, ok := value.Data.(*process.Process)
			if !ok {
				return
			}
			processId, err := strconv.Atoi(value.Id)
			if err != nil {
				return
			}
			data, err := h.m.getProcessData(processId)
			if err != nil {
				return
			}
			if p.Enable != data.process.Enable {
				isEnableChanged = true
			}
			data.process = *p
			if isEnableChanged {
				if data.process.Enable {
					_ = h.m.start(data)
				} else {
					_ = h.m.stop(data)
				}
			}
		}
	}

	if isEnableChanged || key == base.KeyProcessStatus {
		_ = h.m.updateStatistic()
		statistic := new(base.Statistic)
		*statistic = h.m.statistic
		_ = bus.PublishResource(base.KeyStatistic, base.StatusUpdated, "", statistic, nil)
	}
}
