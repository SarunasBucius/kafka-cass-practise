package main

import (
	"context"
	"fmt"
	"log"
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
	db, err := database.SQLiteGormConn()
	if err != nil {
		return err
	}
	defer func() {
		d, _ := db.DB()
		d.Close()
	}()

	k := kcp.New(
		&async.Produce{Producer: prod},
		&database.Gorm{DB: db},
	)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer waitWithTimeout(wg, cancel, time.Second*15)
	if err := startServices(ctx, cancel, k, wg); err != nil {
		return err
	}

	fmt.Println("Hello")
	fmt.Println(version)

	select {
	case <-ctx.Done():
	case <-sig:
	}

	return nil
}

func startServices(ctx context.Context, cancel context.CancelFunc, k *kcp.Kcp, wg *sync.WaitGroup) error {
	{
		cons, err := async.KafkaConsumerConn("inserter")
		if err != nil {
			return err
		}
		wg.Add(1)
		go async.InsertEventsConsumer(ctx, k.InsertVisit, cons, cancel, wg)
	}
	{
		cons, err := async.KafkaConsumerConn("inserter")
		if err != nil {
			return err
		}
		wg.Add(1)
		go async.InsertEventsConsumer(ctx, k.InsertVisit, cons, cancel, wg)
	}
	{
		cons, err := async.KafkaConsumerConn("day", map[string]kafka.ConfigValue{
			"go.events.channel.enable": true,
			"go.events.channel.size":   5,
			"enable.auto.commit":       false,
		})
		if err != nil {
			return err
		}
		wg.Add(1)
		go async.PrintDayConsumer(ctx, k.PrintDay, cons, cancel, wg)
	}

	wg.Add(1)
	go services.ListenHTTP(ctx, services.GinRoutes(k), cancel, wg)

	return nil
}

func waitWithTimeout(wg *sync.WaitGroup, cancel context.CancelFunc, timeout time.Duration) {
	cancel()
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
