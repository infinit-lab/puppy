package main

import (
	_ "github.com/infinit-lab/taiji/src/controller/account"
	_ "github.com/infinit-lab/taiji/src/controller/license"
	_ "github.com/infinit-lab/taiji/src/controller/log"
	_ "github.com/infinit-lab/taiji/src/controller/net"
	_ "github.com/infinit-lab/taiji/src/controller/notification"
	_ "github.com/infinit-lab/taiji/src/controller/performance"
	"github.com/infinit-lab/taiji/src/controller/process"
	_ "github.com/infinit-lab/taiji/src/controller/proxy"
	"github.com/infinit-lab/taiji/src/controller/system"
	_ "github.com/infinit-lab/taiji/src/controller/token"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/signal"
)

var Version string
var CommitId string
var BuildTime string

func main() {
	logutils.Trace("Starting...")
	for i, arg := range os.Args {
		logutils.TraceF("%d. %s", i+1, arg)
	}
	logutils.Trace("Pid is ", os.Getpid())

	isGuard := config.GetBool("process.guard")
	logutils.Trace("IsGuard is ", isGuard)

	if !isGuard {
		go func() {
			v := system.Version{
				Version:   Version,
				CommitId:  CommitId,
				BuildTime: BuildTime,
			}
			system.SetVersion(&v)
			_ = httpserver.ListenAndServe()
		}()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	s := <-c
	logutils.Trace("Got signal: ", s)

	if !isGuard {
		process.Quit()
	}
}
