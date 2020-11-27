// Package kcp provides event production, event consumption
// and data insertion.
package kcp

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

// Consume contains Consumer interface.
type Consume struct {
	Consumer
}

// Consumer interface consumes event.
type Consumer interface {
	ConsumeEvent(Event) error
}

// ConsumeVisit consumes visit event and returns error.
func (c *Consume) ConsumeVisit(event Event) error {
	return c.ConsumeEvent(event)
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
