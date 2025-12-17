package main

import (
	"github.com/cchrisris/go-course/hw2/internal/args"
	"github.com/cchrisris/go-course/hw2/internal/server"
)

func main() {
	port := args.GetPort()

	serv := server.NewServer(port)
	serv.Start()
}
