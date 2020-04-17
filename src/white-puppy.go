package main

import (
	"fmt"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/httpserver"
)

func main() {
	config.Exec()
	fmt.Println("Hello World!")
	_ = httpserver.ListenAndServe()
}
