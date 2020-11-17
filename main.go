package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var version string

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	fmt.Println("Hello")
	fmt.Println(version)

	log.Println(<-sig)
	os.Exit(0)
}
