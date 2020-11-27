// Package kcp provides event production, event consumption
// and data insertion.
package kcp

import "time"

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
type Event time.Time

// Producer produces event.
type Producer interface {
	ProduceEvent(Event) error
}

// ProduceVisit produces visit Event and returns error.
func (k *Kcp) ProduceVisit() error {
	eventTime := time.Now().UTC()
	return k.ProduceEvent(Event(eventTime))
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
