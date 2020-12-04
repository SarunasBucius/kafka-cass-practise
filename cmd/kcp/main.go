package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
	"github.com/SarunasBucius/kafka-cass-practise/platform/async"
	"github.com/SarunasBucius/kafka-cass-practise/platform/database"
	"github.com/SarunasBucius/kafka-cass-practise/platform/services"
)

var version string

func main() {
	err := runApp()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func runApp() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	prod, err := async.KafkaProducerConn()
	if err != nil {
		return err
	}
	defer prod.Close()
	cons, err := async.KafkaConsumerConn()
	if err != nil {
		return err
	}
	defer cons.Close()
	db, err := database.CassConn()
	if err != nil {
		return err
	}
	defer db.Close()

	k := kcp.New(
		&async.Produce{Producer: prod},
		&async.Handle{},
		&database.Insert{Session: db},
	)

	srv := &http.Server{
		Addr:         ":5000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      services.SetRoutes(k),
	}

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go async.ConsumeEvents(ctx, k, cons, cancel, wg)
	wg.Add(1)
	go services.ListenHTTP(ctx, srv, cancel, wg)

	fmt.Println("Hello")
	fmt.Println(version)

	select {
	case <-ctx.Done():
	case <-sig:
	}

	cancel()
	waitWithTimeout(wg, time.Second*15)
	return err
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
	case <-time.After(timeout):
	}
}
