package async

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Produce contains connection to kafka producer.
type Produce struct {
	*kafka.Producer
}

// ProduceEvent produces kcp.Event to kafka.
func (p *Produce) ProduceEvent(event kcp.Event) error {
	topic := "visits"
	b, err := encodeGob(event)
	if err != nil {
		return err
	}
	if err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
		Key:            []byte(event.IP),
	}, nil); err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		return err
	}

	switch e := (<-p.Events()).(type) {
	case *kafka.Message:
		if e.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", e.TopicPartition.Error)
			return e.TopicPartition.Error
		}
	case kafka.Error:
		fmt.Println(e)
		return e
	}

	return nil
}

func encodeGob(event kcp.Event) ([]byte, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b.Bytes(), nil
}
