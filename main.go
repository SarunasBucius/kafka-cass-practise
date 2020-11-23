package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gocql/gocql"
)

var version string

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	prod := kafkaProducerConn()
	defer prod.Close()
	cons := kafkaConsumerConn()
	defer cons.Close()
	db := cassConn()
	defer db.Close()

	checkKafkaConn(cons)

	fmt.Println("Hello")
	fmt.Println(version)

	log.Println(<-sig)
}

func cassConn() *gocql.Session {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}

func kafkaConsumerConn() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_HOST"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	return c
}

func kafkaProducerConn() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_HOST")})
	if err != nil {
		panic(err)
	}
	return p
}

func checkKafkaConn(c *kafka.Consumer) {
	_, err := c.GetMetadata(nil, false, 10000)
	if err != nil {
		panic(err)
	}
}
