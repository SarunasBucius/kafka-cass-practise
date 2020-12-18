// Package async provides event streaming using kafka.
package async

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConsumerConn takes groupID and optional options as param,
// returns connection to kafka consumer or an error.
func KafkaConsumerConn(groupID string, options ...map[string]kafka.ConfigValue) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_HOST"),
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}
	if options != nil {
		for i, val := range options[0] {
			config.SetKey(i, val)
		}
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// KafkaProducerConn returns connection to kafka producer or an error.
func KafkaProducerConn() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_HOST")})
	if err != nil {
		return nil, err
	}
	return p, nil
}
