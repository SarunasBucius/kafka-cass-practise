package async

import (
	"context"
	"fmt"
	"sync"
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
		Value:          []byte(time.Time(event).Format(time.RFC3339)),
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
func ConsumeEvents(ctx context.Context, k *kcp.Kcp, cons *kafka.Consumer, wg *sync.WaitGroup) {
	defer wg.Done()
	err := cons.SubscribeTopics([]string{"visits"}, nil)
	if err != nil {
		fmt.Printf("Subscription failed: %v\n", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := cons.Poll(500)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				event, err := time.Parse(time.RFC3339, string(e.Value))
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = k.HandleVisit(kcp.Event(event))
				if err != nil {
					fmt.Println(err)
				}
			case kafka.Error:
				if e.IsFatal() {
					return
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

// Handle empty struct to use as receiver for HandleEvent
type Handle struct{}

// HandleEvent handles kcp.Event
func (h *Handle) HandleEvent(event kcp.Event) error {
	return nil
}
