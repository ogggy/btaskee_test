package main

import (
	"flag"
	"os"
)

var (
	addr string
)

func init() {
	defer flag.Parse()

	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "http service address")
}

func main() {

	var (
		systemSignal = listenSignal()
	)
	startAPIServer(systemSignal, addr)
	select {
	case <-systemSignal:
		os.Exit(0)
	}
}
