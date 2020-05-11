package process

import (
	"github.com/infinit-lab/taiji/src/model/process"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type guard struct {
	cmd *exec.Cmd
}

func (g *guard) run() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			var err error
			g.cmd = new(exec.Cmd)
			g.cmd.Path = os.Args[0]
			g.cmd.Args = []string{
				os.Args[0],
				"process.guard=true",
				"process.pid=" + strconv.Itoa(os.Getpid()),
			}
			g.cmd.Dir, err = os.Getwd()
			if err != nil {
				logutils.Trace("Failed to Getwd. error: ", err)
				g.cmd = nil
				continue
			}

			if err := g.cmd.Start(); err != nil {
				logutils.Error("Failed to Start. error: ", err)
				g.cmd = nil
				continue
			}
			_ = g.cmd.Wait()
			g.cmd = nil
		}
	}()
}

type slave struct {
	pid int
}

func (s *slave) quit() {
	processList, err := process.GetProcessList()
	if err == nil {
		for _, p := range processList {
			if p.Pid == 0 {
				continue
			}
			pro, err := os.FindProcess(p.Pid)
			if err != nil {
				continue
			}
			err = pro.Kill()
			if err != nil {
				continue
			}
			_, _ = pro.Wait()
			p.Pid = 0
			_ = process.UpdateProcess(p.Id, p, nil)
		}
	}
	os.Exit(0)
}

func (s *slave) run() {
	go func () {
		s.pid = config.GetInt("process.pid")
		logutils.Trace("Get pid is ", s.pid)
		pro, err := os.FindProcess(s.pid)
		if err != nil {
			logutils.Trace("Failed to FindProcess. error: ", err)
			s.quit()
		}
		logutils.Trace("Wait process ", s.pid)
		_, _ = pro.Wait()
		logutils.TraceF("Process %d quit.", s.pid)
		s.quit()
	}()
}
