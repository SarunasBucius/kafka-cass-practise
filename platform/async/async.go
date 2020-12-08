package async

import (
	"context"
	"fmt"
	"os"
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
		Key:            []byte{},
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

// InsertEventsConsumer inserts events from kafka consumer.
func InsertEventsConsumer(ctx context.Context, k *kcp.Kcp, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := cons.SubscribeTopics([]string{"visits"}, nil); err != nil {
		fmt.Printf("Subscription failed: %v\n", err)
		cancel()
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
				fmt.Println(cons.String())
				event, err := time.Parse(time.RFC3339, string(e.Value))
				if err != nil {
					fmt.Println(err)
					continue
				}
				if err := k.InsertVisit(kcp.Event(event)); err != nil {
					fmt.Println(err)
				}
			case kafka.Error:
				if e.IsFatal() {
					cancel()
					return
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

// PrintDayConsumer prints day from consumed events.
func PrintDayConsumer(ctx context.Context, k *kcp.Kcp, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := cons.SubscribeTopics([]string{"visits"}, nil); err != nil {
		cancel()
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-cons.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Println(cons.String())
				event, err := time.Parse(time.RFC3339, string(e.Value))
				if err != nil {
					fmt.Println(err)
					continue
				}
				wg.Add(1)
				go func(event time.Time) {
					defer wg.Done()
					k.PrintDay(kcp.Event(event))
				}(event)
			case kafka.Error:
				if e.IsFatal() {
					cancel()
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

// KafkaConsumerConn takes groupID and optional options as param, returns connection to kafka consumer or an error.
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

// KafkaProducerConn returns connection to kafka producer or an error
func KafkaProducerConn() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_HOST")})
	if err != nil {
		return nil, err
	}
	return p, nil
}
