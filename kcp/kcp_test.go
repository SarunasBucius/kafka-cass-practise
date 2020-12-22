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

func TestFormatTime(t *testing.T) {
	type test struct {
		filter map[string]string
		key    string
		want   time.Time
		err    error
	}

	y, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatal(err)
	}
	ymd, err := time.Parse("2006-01-02", "2020-05-05")
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]test{
		"no value": {
			filter: map[string]string{},
			key:    "date",
			want:   time.Time{},
			err:    nil,
		},
		"year": {
			filter: map[string]string{"date": "2020"},
			key:    "date",
			want:   y,
			err:    nil,
		},
		"year month day": {
			filter: map[string]string{"date": "2020-05-05"},
			key:    "date",
			want:   ymd,
			err:    nil,
		},
		"invalid year": {
			filter: map[string]string{"date": "abc"},
			key:    "date",
			want:   time.Time{},
			err:    ErrInvalidFilter,
		},
	}

	for name, tt := range tests {
		got, err := formatTime(tt.filter, tt.key)
		if got != tt.want {
			t.Errorf("%s: expected: %v, got: %v", name, tt.want, got)
		}
		if err != tt.err {
			t.Errorf("%s: expected: %v, got: %v", name, tt.err, err)
		}
	}
}

func TestIsValidDay(t *testing.T) {
	type test struct {
		filter map[string]string
		want   string
		err    error
	}

	tests := map[string]test{
		"no value": {
			filter: map[string]string{},
			want:   "",
			err:    nil,
		},
		"valid day": {
			filter: map[string]string{"day": "Monday"},
			want:   "Monday",
			err:    nil,
		},
		"invalid day": {
			filter: map[string]string{"day": "Mday"},
			want:   "",
			err:    ErrInvalidFilter,
		},
	}

	for name, tt := range tests {
		got, err := isValidDay(tt.filter)
		if got != tt.want {
			t.Errorf("%s: expected: %v, got: %v", name, tt.want, got)
		}
		if err != tt.err {
			t.Errorf("%s: expected: %v, got: %v", name, tt.err, err)
		}
	}
}
