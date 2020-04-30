package main

import (
	_ "github.com/infinit-lab/puppy/src/controller/account"
	_ "github.com/infinit-lab/puppy/src/controller/notification"
	_ "github.com/infinit-lab/puppy/src/controller/token"
	"github.com/infinit-lab/yolanda/httpserver"
)

func main() {
	_ = httpserver.ListenAndServe()
}
