package kcp

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestProduceVisit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProdEvent := NewMockProducer(mockCtrl)
	k := New(mockProdEvent, nil)

	param := "ip"
	mockProdEvent.EXPECT().
		ProduceEvent(approxTime{dev: time.Second, ip: param}).
		Return(nil).
		Times(1)

	k.ProduceVisit(param)
}

type approxTime struct {
	dev time.Duration
	ip  string
}

func (a approxTime) Matches(event interface{}) bool {
	ev, ok := event.(Event)
	if !ok {
		return false
	}

	now := time.Now().UTC()
	if now.Sub(ev.VisitedAt) > a.dev || now.Sub(ev.VisitedAt) < 0 {
		return false
	}

	if ev.VisitedAt.Weekday().String() != now.Weekday().String() {
		return false
	}

	if a.ip != ev.IP {
		return false
	}
	return true
}

func (a approxTime) String() string {
	now := time.Now().UTC()
	event := Event{VisitedAt: now, Day: now.Weekday().String(), IP: a.ip}
	return fmt.Sprintf("%v, with deviation of %v", event, a.dev)
}

func TestInsertVisit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDb := NewMockDbConnector(mockCtrl)
	k := New(nil, mockDb)

	event := Event{}
	mockDb.EXPECT().InsertEvent(event).Return(nil)

	k.InsertVisit(event)
}
