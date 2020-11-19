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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	fmt.Println("Hello")
	fmt.Println(version)

	log.Println(<-sig)
}
