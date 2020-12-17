package database

import (
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Db contains connection to cassandra db.
type Db struct {
	*gocql.Session
}

// InsertEvent inserts kcp.Event into cassandra db
func (db *Db) InsertEvent(e kcp.Event) error {
	return db.Query(
		"INSERT INTO kcp.visits (ip, visited_at, day) VALUES (?, ?, ?)",
		e.IP,
		e.VisitedAt,
		e.Day).Exec()
}

// GetVisits get visits grouped by ip
func (db *Db) GetVisits() (kcp.VisitsByIP, error) {
	iter := db.Query("SELECT ip, visited_at FROM kcp.visits").Iter()

	var ip string
	var t time.Time
	visits := make(kcp.VisitsByIP)
	for iter.Scan(&ip, &t) {
		visits[ip] = append(visits[ip], t)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return visits, nil
}

// CassConn returns connection to cassandra db or an error
func CassConn() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Timeout = time.Second * 2
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
