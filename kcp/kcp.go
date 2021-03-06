// Package kcp provides visits tracking.
//
// Supported features:
//  * Produce event of visit
//  * Insert event of visit to storage
//  * Print week day of visit
//  * Get all events from storage
//   * Filters visits by time greater than (gt), less than (lt), day of the week (day)
//  * Get events by same ip from storage
//   * Filters visits by time greater than (gt), less than (lt), day of the week (day)
package kcp

//go:generate mockgen -destination=kcp_mock.go -package=kcp github.com/SarunasBucius/kafka-cass-practise/kcp Producer,DbConnector

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Kcp contains Producer and DbConnector.
type Kcp struct {
	Producer
	DbConnector
}

// New takes Producer, DbConnector and returns Kcp instance.
func New(p Producer, i DbConnector) *Kcp {
	return &Kcp{Producer: p, DbConnector: i}
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
	day := now.Weekday().String()
	event := Event{VisitedAt: now, IP: ip, Day: day}
	return k.ProduceEvent(event)
}

// DbConnector interface contains methods concerned with database.
type DbConnector interface {
	InsertEvent(Event) error
	GetVisits() (VisitsByIP, error)
	GetVisitsByIP(ip, day string, gt, lt time.Time) (VisitsByIP, error)
}

// InsertVisit inserts visit Event and returns error.
func (k *Kcp) InsertVisit(event Event) error {
	return k.InsertEvent(event)
}

// VisitsByIP contains ip and slice of visit times.
type VisitsByIP map[string][]time.Time

// ErrInvalidFilter is returned if filter parameter is invalid.
var ErrInvalidFilter = errors.New("invalid filter parameter")

// GetVisits get visits grouped by ip.
func (k *Kcp) GetVisits(filter map[string]string) (VisitsByIP, error) {
	// check if filter for greater than is passed and get valid time.Time value
	gt, err := formatTime(filter, "gt")
	if err != nil {
		return nil, err
	}

	// check if filter for less than is passed and get valid time.Time value
	lt, err := formatTime(filter, "lt")
	if err != nil {
		return nil, err
	}

	// check if filter for day is passed and is valid
	day, err := isValidDay(filter)
	if err != nil {
		return nil, err
	}

	// get visits from db
	visits, err := k.DbConnector.GetVisits()
	if err != nil {
		return nil, err
	}

	// filter data
	filteredVisits := make(VisitsByIP)
	for i, visitsByIP := range visits {
		var filteredTime []time.Time
		for _, visit := range visitsByIP {
			if !gt.IsZero() && visit.Before(gt) {
				continue
			}
			if !lt.IsZero() && visit.After(lt) {
				continue
			}
			if day != "" && day != visit.Weekday().String() {
				continue
			}
			filteredTime = append(filteredTime, visit)
		}
		if filteredTime == nil {
			continue
		}
		filteredVisits[i] = filteredTime
	}
	return filteredVisits, nil
}

func isValidDay(filter map[string]string) (string, error) {
	day := filter["day"]
	if day == "" {
		return "", nil
	}
	var empty struct{}
	days := map[string]struct{}{
		"Monday":    empty,
		"Tuesday":   empty,
		"Wednesday": empty,
		"Thursday":  empty,
		"Friday":    empty,
		"Saturday":  empty,
		"Sunday":    empty,
	}
	if _, ok := days[day]; ok {
		return day, nil
	}
	return "", ErrInvalidFilter
}

func formatTime(filter map[string]string, key string) (time.Time, error) {
	// check if value is passed
	unf := filter[key]
	if unf == "" {
		return time.Time{}, nil
	}

	// add month and day if missing
	for i := len(strings.Split(unf, "-")); i < 3; i++ {
		unf += "-01"
	}

	f, err := time.Parse("2006-01-02", unf)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, ErrInvalidFilter
	}
	return f, nil
}

// GetVisitsByIP gets visits from provided ip.
func (k *Kcp) GetVisitsByIP(ip string, filter map[string]string) (VisitsByIP, error) {
	// check if filter for greater than is passed and get valid time.Time value
	gt, err := formatTime(filter, "gt")
	if err != nil {
		return nil, err
	}

	// check if filter for less than is passed and get valid time.Time value
	lt, err := formatTime(filter, "lt")
	if err != nil {
		return nil, err
	}

	// check if filter for day is passed and is valid
	day, err := isValidDay(filter)
	if err != nil {
		return nil, err
	}
	return k.DbConnector.GetVisitsByIP(ip, day, gt, lt)
}

// PrintDay prints day of the week of event.
func (*Kcp) PrintDay(event Event) {
	fmt.Println(event.Day)
}
