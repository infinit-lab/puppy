package main

import (
	_ "github.com/infinit-lab/taiji/src/controller/account"
	_ "github.com/infinit-lab/taiji/src/controller/notification"
	_ "github.com/infinit-lab/taiji/src/controller/performance"
	_ "github.com/infinit-lab/taiji/src/controller/token"
	"github.com/infinit-lab/yolanda/httpserver"
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/signal"
)

func main() {
	go func() {
		_ = httpserver.ListenAndServe()
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	s := <- c
	logutils.Trace("Got signal: ", s)
}
