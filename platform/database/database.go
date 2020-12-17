package database

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Db contains connection to cassandra db.
type Db struct {
	*gocql.Session
}

// InsertEvent inserts kcp.Event into cassandra db
func (i *Db) InsertEvent(e kcp.Event) error {
	return i.Query(
		"INSERT INTO kcp.visits (id, visited_at) VALUES (?, ?)",
		gocql.TimeUUID(),
		e.VisitedAt).Exec()
}

// CassConn returns connection to cassandra db or an error
func CassConn() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if err := initDb(session); err != nil {
		return nil, err
	}
	return session, nil
}

func initDb(s *gocql.Session) error {
	fmt.Println("Init database")
	if err := s.Query(`DROP KEYSPACE IF EXISTS kcp`).Exec(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := s.Query(`
	CREATE  KEYSPACE IF NOT EXISTS kcp 
	WITH REPLICATION = { 
		'class' : 'SimpleStrategy',
		'replication_factor' : 1 }`,
	).Exec(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := s.Query(`
	CREATE TABLE kcp.visits(
		ip text,
		visited_at timestamp,
		day text,
		PRIMARY KEY (ip, visited_at))`,
	).Exec(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := s.Query(`
	CREATE INDEX IF NOT EXISTS ON kcp.visits (day)`,
	).Exec(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
