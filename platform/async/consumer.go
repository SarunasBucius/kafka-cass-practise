package async

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// InsertVisit describes method to insert visit.
type InsertVisit func(kcp.Event) error

// InsertEventsConsumer inserts events from kafka consumer.
func InsertEventsConsumer(ctx context.Context, insertVisit InsertVisit, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
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
				if err := insertVisit(event); err != nil {
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

// PrintDay describes method to print day.
type PrintDay func(kcp.Event)

// PrintDayConsumer prints day from consumed events.
func PrintDayConsumer(ctx context.Context, printDay PrintDay, cons *kafka.Consumer, cancel context.CancelFunc, wg *sync.WaitGroup) {
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
		case <-time.After(time.Second * 5):
			offsetDif, _ = commitOffset(offsetDif, 1, cons)
		case <-done:
			offsetDif++
		case ev := <-cons.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				wg.Add(1)
				go asyncPrintDay(e.Value, printDay, wg, done)
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

func asyncPrintDay(value []byte, printDay PrintDay, wg *sync.WaitGroup, done chan<- struct{}) {
	defer wg.Done()
	defer func() {
		done <- struct{}{}
	}()
	event, err := decodeGob(value)
	if err != nil {
		fmt.Println(err)
		return
	}
	printDay(event)
}

func decodeGob(data []byte) (kcp.Event, error) {
	var event kcp.Event
	if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&event); err != nil {
		fmt.Println(err)
		return kcp.Event{}, err
	}
	return event, nil
}
