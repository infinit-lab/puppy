package search

import (
	"github.com/infinit-lab/yolanda/logutils"
	"os"
	"os/signal"
	"testing"
)

func TestInit(t *testing.T) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	s := <- c
	logutils.Trace("Got signal: ", s)
}
