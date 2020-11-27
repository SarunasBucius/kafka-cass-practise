package async

import "github.com/confluentinc/confluent-kafka-go/kafka"

// Produce contains connection to kafka producer
type Produce struct {
	*kafka.Producer
}

// ProduceEvent produces event to kafka
func (p Produce) ProduceEvent() error {

	return nil
}
