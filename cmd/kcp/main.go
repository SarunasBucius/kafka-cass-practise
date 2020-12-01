package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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

	err = checkKafkaConn(cons)
	if err != nil {
		return err
	}

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

	go async.ConsumeEvents(k, cons)
	go listenHTTP(srv)

	fmt.Println("Hello")
	fmt.Println(version)

	log.Println(<-sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	return nil
}

func listenHTTP(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
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

func checkKafkaConn(c *kafka.Consumer) error {
	fmt.Println("Checking initial connection")
	md, err := c.GetMetadata(nil, false, 10000)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Originating broker: %v\n", md.OriginatingBroker)
	fmt.Printf("Brokers: %v\n", md.Brokers)
	h, err := net.LookupHost(md.Brokers[0].Host)
	if err != nil {
		return err
	} else if len(h) < 1 {
		return errors.New("host address not found")
	}
	fmt.Println(h)
	if _, ok := md.Topics["testTopic"]; !ok {
		err = insertTopic(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertTopic(c *kafka.Consumer) error {
	fmt.Println("Creating new topic")
	a, err := kafka.NewAdminClientFromConsumer(c)
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		return err
	}

	defer a.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = a.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             "testTopic",
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(time.Second*15))
	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		return err
	}
	return nil

}
