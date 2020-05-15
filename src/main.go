package main

import (
	_ "github.com/infinit-lab/taiji/src/controller/account"
	_ "github.com/infinit-lab/taiji/src/controller/notification"
	_ "github.com/infinit-lab/taiji/src/controller/performance"
	"github.com/infinit-lab/taiji/src/controller/process"
	_ "github.com/infinit-lab/taiji/src/controller/token"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/signal"
)

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
