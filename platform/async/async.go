package async

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Produce contains connection to kafka producer
type Produce struct {
	*kafka.Producer
}

// ProduceEvent produces kcp.Event to kafka
func (p *Produce) ProduceEvent(event kcp.Event) error {
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
