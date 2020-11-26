// Package kcp provides event production, event consumption
// and data insertion
package kcp

// Producer produces events.
type Producer interface {
	ProduceEvent() error
}

// ProduceVisits produces visits and returns error.
func ProduceVisits(p Producer) error {
	return p.ProduceEvent()
}
