package async

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Produce contains connection to kafka producer
type Produce struct {
	*kafka.Producer
}

// ProduceEvent produces kcp.Event to kafka
func (p *Produce) ProduceEvent(event kcp.Event) error {
	topic := "visits"
	if err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(time.Time(event).String()),
	}, nil); err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		return err
	}

	e := <-p.Events()
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}
	return nil
}

// ConsumeEvents consumes events from kafka
func ConsumeEvents(k *kcp.Kcp, cons *kafka.Consumer) {

}

// Handle empty struct to use as receiver for HandleEvent
type Handle struct{}

// HandleEvent handles kcp.Event
func (h *Handle) HandleEvent(event kcp.Event) error {
	return nil
}
