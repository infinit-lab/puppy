package main

import (
	_ "github.com/infinit-lab/taiji/src/controller/account"
	_ "github.com/infinit-lab/taiji/src/controller/notification"
	_ "github.com/infinit-lab/taiji/src/controller/performance"
	_ "github.com/infinit-lab/taiji/src/controller/token"
	"github.com/infinit-lab/yolanda/httpserver"
)

func main() {
	_ = httpserver.ListenAndServe()
}
