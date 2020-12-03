package database

import (
	"time"

	"github.com/gocql/gocql"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Insert contains connection to cassandra db
type Insert struct {
	*gocql.Session
}

// InsertEvent inserts kcp.Event into cassandra db
func (i *Insert) InsertEvent(e kcp.Event) error {
	return i.Query(
		"INSERT INTO kcp.visits (id, visited_at) VALUES (?, ?)",
		gocql.TimeUUID(),
		time.Time(e)).Exec()
}
