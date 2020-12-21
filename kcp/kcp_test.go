package kcp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
	"github.com/SarunasBucius/kafka-cass-practise/mocks"
)

func TestProduceVisit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProdEvent := mocks.NewMockProducer(mockCtrl)
	k := kcp.New(mockProdEvent, nil)

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
	ev, ok := event.(kcp.Event)
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
	event := kcp.Event{VisitedAt: now, Day: now.Weekday().String(), IP: a.ip}
	return fmt.Sprintf("%v, with deviation of %v", event, a.dev)
}
