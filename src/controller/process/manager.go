package process

import (
	"errors"
	"github.com/infinit-lab/taiji/src/model/base"
	"github.com/infinit-lab/taiji/src/model/process"
	"github.com/infinit-lab/yolanda/bus"
	"github.com/infinit-lab/yolanda/logutils"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type processData struct {
	process process.Process
	cmd     exec.Cmd
	isStart bool
	mutex   sync.Mutex
}

type manager struct {
	processList []*processData
	statistic   base.Statistic
	mutex       sync.Mutex
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
		data.cmd.Path = p.Path
		data.cmd.Args = strings.Split(p.Config, " ")
		data.cmd.Dir = p.Dir
		if data.process.Enable {
			_ = m.start(data)
		} else {
			data.isStart = false
			updateProcessStatus(data, false, 0)
		}
	}
}

func updateProcessStatus(p *processData, started bool, pid int) {
	status := process.Status{
		ProcessId: p.process.Id,
		Type:      base.StatusTypeStarted,
	}
	if started {
		status.Value = "1"
		p.process.Pid = pid
	} else {
		status.Value = "0"
		p.process.Pid = 0
	}
	_ = process.UpdateStatus(&status, nil)
	_ = process.UpdateProcess(p.process.Id, &p.process, nil)
}

func (m *manager) start(p *processData) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isStart {
		return nil
	}
	p.isStart = true
	go func() {
		for p.isStart {
			time.Sleep(100 * time.Millisecond)
			if err := p.cmd.Start(); err != nil {
				logutils.ErrorF("Failed to start %s. error: %v", p.process.Name, err)
				updateProcessStatus(p, false, 0)
				continue
			}
			updateProcessStatus(p, true, p.cmd.Process.Pid)

			if err := p.cmd.Wait(); err != nil {
				logutils.WarningF("%s quit. error: %v", p.process.Name, err)
			} else {
				logutils.WarningF("%s quit.", p.process.Name)
			}
			updateProcessStatus(p, false, 0)
		}
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
	return p.cmd.Process.Kill()
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

type processHandler struct {
	m *manager
}

func (h *processHandler) Handle(key int, value *bus.Resource) {
	if key == base.KeyProcess || key == base.KeyProcessEnable {
		if value.Status == base.StatusUpdated {
			p, ok := value.Data.(*process.Process)
			if !ok {
				return
			}
			data, err := h.m.getProcessData(key)
			if err != nil {
				return
			}
			isEnableChanged := false
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
	_ = h.m.updateStatistic()
	_ = bus.PublishResource(base.KeyStatistic, base.StatusUpdated, "", &h.m.statistic, nil)
}
