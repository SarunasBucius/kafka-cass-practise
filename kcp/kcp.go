// Package kcp provides event production, event consumption
// and data insertion
package kcp

// ProdConf contains Produces interface
type ProdConf struct {
	Producer
}

// Producer produces events.
type Producer interface {
	ProduceEvent() error
}

// ProduceVisits produces visits and returns error.
func (p *ProdConf) ProduceVisits() error {
	return p.ProduceEvent()
}
