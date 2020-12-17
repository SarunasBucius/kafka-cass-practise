// Package kcp provides event production, event consumption,
// data insertion and data view.
package kcp

import (
	"fmt"
	"time"
)

// Kcp contains Producer, Handler and DbConnector.
type Kcp struct {
	Producer
	Handler
	DbConnector
}

// New takes Producer, Handler, DbConnector and returns Kcp instance.
func New(p Producer, h Handler, i DbConnector) *Kcp {
	return &Kcp{Producer: p, Handler: h, DbConnector: i}
}

// Event represents event created by ProduceVisit.
type Event struct {
	VisitedAt time.Time
	IP        string
	Day       string
}

// Producer produces event.
type Producer interface {
	ProduceEvent(Event) error
}

// ProduceVisit takes ip as param, produces visit Event and returns error.
func (k *Kcp) ProduceVisit(ip string) error {
	now := time.Now().UTC()
	day := time.Time(now).Weekday().String()
	event := Event{VisitedAt: now, IP: ip, Day: day}
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

// DbConnector interface contains methods concerned with database.
type DbConnector interface {
	InsertEvent(Event) error
	GetVisits() (VisitsByIP, error)
}

// InsertVisit inserts visit Event and returns error.
func (k *Kcp) InsertVisit(event Event) error {
	return k.InsertEvent(event)
}

// VisitsByIP contains ip and slice of visit times
type VisitsByIP map[string][]time.Time

// GetVisits get visits grouped by ip
func (k *Kcp) GetVisits() (VisitsByIP, error) {
	return k.DbConnector.GetVisits()
}

// PrintDay prints day of the week of event
func (*Kcp) PrintDay(event Event) {
	fmt.Println(event.Day)
}
