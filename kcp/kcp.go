// Package kcp provides event production, event consumption
// and data insertion.
package kcp

import (
	"fmt"
	"time"
)

// Kcp contains Producer, Handler and Insert
type Kcp struct {
	Producer
	Handler
	Inserter
}

// New takes Producer, Handler, Inserter and returns Kcp instance
func New(p Producer, h Handler, i Inserter) *Kcp {
	return &Kcp{Producer: p, Handler: h, Inserter: i}
}

// Event represents event created by ProduceVisit.
type Event struct {
	VisitedAt time.Time
	IP        string
}

// Producer produces event.
type Producer interface {
	ProduceEvent(Event) error
}

// ProduceVisit takes ip as param, produces visit Event and returns error.
func (k *Kcp) ProduceVisit(ip string) error {
	event := Event{VisitedAt: time.Now().UTC(), IP: ip}
	return k.ProduceEvent(event)
}

// Handler interface handles event.
type Handler interface {
	HandleEvent(Event) error
}

// HandleVisit handles visit Event and returns error.
func (k *Kcp) HandleVisit(event Event) error {
	return k.HandleEvent(event)
}

// Inserter interface inserts event.
type Inserter interface {
	InsertEvent(Event) error
}

// InsertVisit inserts visit Event and returns error.
func (k *Kcp) InsertVisit(event Event) error {
	return k.InsertEvent(event)
}

// PrintDay prints day of the week of event
func (*Kcp) PrintDay(event Event) {
	fmt.Println(time.Time(event.VisitedAt).Weekday())
}
