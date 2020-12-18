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

// GetVisitsByIP get filtered visits by ip.
func (db *Db) GetVisitsByIP(ip, day string, gt, lt time.Time) (kcp.VisitsByIP, error) {
	var params []interface{}
	q := fmt.Sprint(`
	SELECT visited_at 
	FROM kcp.visits 
	WHERE ip=?
	`)
	params = append(params, ip)

	if !gt.IsZero() {
		q = fmt.Sprintf("%v AND visited_at > ?", q)
		params = append(params, gt)
	}

	if !lt.IsZero() {
		q = fmt.Sprintf("%v AND visited_at < ?", q)
		params = append(params, lt)
	}

	if day != "" {
		q = fmt.Sprintf("%v AND day = ?", q)
		params = append(params, day)
	}
	iter := db.Query(q, params...).Iter()

	var t time.Time
	visits := make(kcp.VisitsByIP)
	for iter.Scan(&t) {
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

	return initialData(s)
}

func initialData(s *gocql.Session) error {
	for i := 0; i < 50; i++ {
		if err := s.Query(
			"INSERT INTO kcp.visits (ip, visited_at, day) VALUES (?, ?, ?)",
			"172.19.0."+fmt.Sprint(i%5),
			time.Now().UTC().AddDate(0, i%5, i),
			time.Now().UTC().AddDate(0, i%5, i).Weekday().String(),
		).Exec(); err != nil {
			return err
		}
	}
	return nil
}
