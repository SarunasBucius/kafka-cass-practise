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

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go async.ConsumeEvents(ctx, k, cons, cancel, wg)
	wg.Add(1)
	go listenHTTP(ctx, srv, cancel, wg)

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

func listenHTTP(ctx context.Context, srv *http.Server, cancel context.CancelFunc, wg *sync.WaitGroup) {
	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(context.Background())
		}
	}()
	defer wg.Done()
	if err := srv.ListenAndServe(); err != nil {
		cancel()
	}
}

func cassConn() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if err := initDb(session); err != nil {
		return nil, err
	}
	return session, nil
}

func initDb(s *gocql.Session) error {
	fmt.Println("Init database")
	err := s.Query(`DROP KEYSPACE IF EXISTS kcp`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.Query(`CREATE  KEYSPACE IF NOT EXISTS kcp 
			WITH REPLICATION = { 
	   		'class' : 'SimpleStrategy',
			'replication_factor' : 1 }`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.Query(`CREATE TABLE kcp.visits(
	id UUID PRIMARY KEY,
	visited_at timestamp)`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
