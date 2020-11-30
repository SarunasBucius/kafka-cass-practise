package database

import (
	"github.com/gocql/gocql"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Insert contains connection to cassandra db
type Insert struct {
	*gocql.Session
}

// InsertEvent inserts kcp.Event into cassandra db
func (i *Insert) InsertEvent(e kcp.Event) error {
	return nil
}
