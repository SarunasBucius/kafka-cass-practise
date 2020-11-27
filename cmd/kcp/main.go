package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
	"github.com/SarunasBucius/kafka-cass-practise/platform/async"
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

	go listenHTTP(prod)

	fmt.Println("Hello")
	fmt.Println(version)

	log.Println(<-sig)
	return nil
}

func listenHTTP(prod *kafka.Producer) {
	p := kcp.New(async.Produce{Producer: prod}, nil, nil)
	r := services.SetRoutes(p)

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		panic(err)
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
	_, err := c.GetMetadata(nil, false, 10000)
	if err != nil {
		return err
	}
	return nil
}