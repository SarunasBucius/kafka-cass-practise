package database

import (
	"fmt"
	"os"
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
	err := s.Query(`DROP KEYSPACE IF EXISTS kcp`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.Query(`CREATE  KEYSPACE IF NOT EXISTS kcp 
			WITH REPLICATION = { 
	   		'class' : 'SimpleStrategy',
			'replication_factor' : 1 }`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.Query(`CREATE TABLE kcp.visits(
	id UUID PRIMARY KEY,
	visited_at timestamp)`).Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
