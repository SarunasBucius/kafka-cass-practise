package async

import (
	"bytes"
	"context"
	"encoding/gob"
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

// VisitInserter contains method to insert visit.
type VisitInserter interface {
	InsertVisit(kcp.Event) error
}

// InsertEventsConsumer inserts events from kafka consumer.
func InsertEventsConsumer(ctx context.Context, v VisitInserter, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	defer cons.Close()
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
				event, err := decodeGob(e.Value)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if err := v.InsertVisit(event); err != nil {
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

// DayPrinter contains method to print day
type DayPrinter interface {
	PrintDay(kcp.Event)
}

// PrintDayConsumer prints day from consumed events.
func PrintDayConsumer(ctx context.Context, d DayPrinter, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	defer cons.Close()
	defer cons.Commit()
	if err := cons.SubscribeTopics([]string{"visits"}, nil); err != nil {
		cancel()
		return
	}

	offsetDif := 0
	done := make(chan struct{}, 5)
	for {
		offsetDif, _ = commitOffset(offsetDif, 5, cons)
		select {
		case <-ctx.Done():
			return
		// Or could just leave default auto commit every X seconds
		case <-time.After(time.Second * 5):
			offsetDif, _ = commitOffset(offsetDif, 1, cons)
		case <-done:
			offsetDif++
		case ev := <-cons.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				wg.Add(1)
				go asyncPrintDay(e.Value, d, wg, done)
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

func commitOffset(offset, minOffset int, cons *kafka.Consumer) (int, error) {
	if offset >= minOffset {
		if _, err := cons.Commit(); err != nil {
			fmt.Println(err)
			return offset, err
		}
		return 0, nil
	}
	return offset, nil
}

func asyncPrintDay(value []byte, d DayPrinter, wg *sync.WaitGroup, done chan<- struct{}) {
	defer wg.Done()
	defer func() {
		done <- struct{}{}
	}()
	event, err := decodeGob(value)
	if err != nil {
		fmt.Println(err)
		return
	}
	d.PrintDay(event)
}

func encodeGob(event kcp.Event) ([]byte, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeGob(data []byte) (kcp.Event, error) {
	var event kcp.Event
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&event); err != nil {
		fmt.Println(err)
		return kcp.Event{}, err
	}
	return event, nil
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
