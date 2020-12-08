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

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
	"github.com/SarunasBucius/kafka-cass-practise/platform/async"
	"github.com/SarunasBucius/kafka-cass-practise/platform/database"
	"github.com/SarunasBucius/kafka-cass-practise/platform/services"
)

var version string

func main() {
	if err := runApp(); err != nil {
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
	cins1, err := async.KafkaConsumerConn("inserter")
	if err != nil {
		return err
	}
	defer cins1.Close()
	cins2, err := async.KafkaConsumerConn("inserter")
	if err != nil {
		return err
	}
	defer cins2.Close()
	cday, err := async.KafkaConsumerConn("day", map[string]kafka.ConfigValue{
		"go.events.channel.enable": true,
		"go.events.channel.size":   5,
	})
	if err != nil {
		return err
	}
	defer cday.Close()

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
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      services.SetRoutes(k),
	}

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go async.InsertEventsConsumer(ctx, k, cins1, cancel, wg)
	wg.Add(1)
	go async.InsertEventsConsumer(ctx, k, cins2, cancel, wg)
	wg.Add(1)
	go async.PrintDayConsumer(ctx, k, cday, cancel, wg)

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
