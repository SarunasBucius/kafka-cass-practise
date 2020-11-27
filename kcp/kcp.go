// Package kcp provides event production, event consumption
// and data insertion.
package kcp

// Kcp contains Producer, Handler and Insert
type Kcp struct {
	Producer
	Handler
	Inserter
}

// New takes Producer, Handler, Inserter and returns Kcp instance and an error
func New(p Producer, h Handler, i Inserter) *Kcp {
	return &Kcp{Producer: p, Handler: h, Inserter: i}
}

// Producer produces event.
type Producer interface {
	ProduceEvent() error
}

// ProduceVisit produces visit and returns error.
func (k *Kcp) ProduceVisit() error {
	return k.ProduceEvent()
}

// Event represets string value of event.
type Event string

// Handler interface handles event.
type Handler interface {
	HandleEvent(Event) error
}

// HandleVisit handles visit event and returns error.
func (k *Kcp) HandleVisit(event Event) error {
	return k.HandleEvent(event)
}

// Inserter interface inserts event.
type Inserter interface {
	InsertEvent(Event) error
}

// InsertVisit inserts visit event and returns error.
func (k *Kcp) InsertVisit(event Event) error {
	return k.InsertEvent(event)
}
