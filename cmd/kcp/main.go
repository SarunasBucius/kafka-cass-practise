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
	"github.com/gocql/gocql"

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

	prod, err := kafkaProducerConn()
	if err != nil {
		return err
	}
	defer prod.Close()
	cons, err := kafkaConsumerConn()
	if err != nil {
		return err
	}
	defer cons.Close()
	db, err := cassConn()
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

	errc := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go async.ConsumeEvents(ctx, k, cons, wg)
	go startListen(ctx, srv, errc, wg)

	fmt.Println("Hello")
	fmt.Println(version)

	select {
	case err = <-errc:
	case <-sig:
	}

	cancel()
	wg.Wait()
	return err
}

func startListen(ctx context.Context, srv *http.Server, errc chan<- error, wg *sync.WaitGroup) {
	err := listenHTTP(ctx, srv, wg)
	if err != nil {
		errc <- err
	}
}

func listenHTTP(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) error {
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	defer wg.Done()
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func cassConn() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func kafkaConsumerConn() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_HOST"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func kafkaProducerConn() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_HOST")})
	if err != nil {
		return nil, err
	}
	return p, nil
}
