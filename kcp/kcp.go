// Package kcp provides event production, event consumption
// and data insertion.
package kcp

import "errors"

// Kcp contains Producer, Handler and Insert
type Kcp struct {
	*Produce
	*Handle
	*Insert
}

// New takes Producer, Handler, Inserter and returns Kcp instance and an error
func New(p Producer, h Handler, i Inserter) (*Kcp, error) {
	if p == nil || h == nil || i == nil {
		return nil, errors.New("missing implementations")
	}
	return &Kcp{
		&Produce{p},
		&Handle{h},
		&Insert{i},
	}, nil
}

// Produce contains Producer interface.
type Produce struct {
	Producer
}

// Producer produces event.
type Producer interface {
	ProduceEvent() error
}

// NewProduce takes Producer as param and creates new Produce instance
func NewProduce(p Producer) *Produce {
	return &Produce{Producer: p}
}

// ProduceVisit produces visit and returns error.
func (p *Produce) ProduceVisit() error {
	return p.ProduceEvent()
}

// Event represets string value of event.
type Event string

// Handle contains Handler interface.
type Handle struct {
	Handler
}

// Handler interface handles event.
type Handler interface {
	HandleEvent(Event) error
}

// HandleVisit handles visit event and returns error.
func (c *Handle) HandleVisit(event Event) error {
	return c.HandleEvent(event)
}

// Insert contains Inserter interface.
type Insert struct {
	Inserter
}

// Inserter interface inserts event.
type Inserter interface {
	InsertEvent(Event) error
}

// InsertVisit inserts visit event and returns error.
func (i *Insert) InsertVisit(event Event) error {
	return i.InsertEvent(event)
}
