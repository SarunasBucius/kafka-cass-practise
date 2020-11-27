package async

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Produce contains connection to kafka producer
type Produce struct {
	*kafka.Producer
}

// ProduceEvent produces event to kafka
func (p Produce) ProduceEvent(e kcp.Event) error {
	return nil
}
