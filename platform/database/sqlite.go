package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// Register to sql package.
	_ "github.com/mattn/go-sqlite3"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// SQLite contains connection to SQLite db.
type SQLite struct {
	*sql.DB
}

// InsertEvent inserts kcp.Event into db.
func (db *SQLite) InsertEvent(e kcp.Event) error {
	_, err := db.Exec(
		"INSERT INTO visits (ip, visited_at, day) VALUES (?, ?, ?)",
		e.IP,
		e.VisitedAt,
		e.Day)
	return err
}

// GetVisits get visits grouped by ip.
func (db *SQLite) GetVisits() (kcp.VisitsByIP, error) {
	rows, err := db.Query("SELECT ip, visited_at FROM visits")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	var ip string
	var t time.Time
	visits := make(kcp.VisitsByIP)
	for rows.Next() {
		if err := rows.Scan(&ip, &t); err != nil {
			fmt.Println(err)
			return nil, err
		}
		visits[ip] = append(visits[ip], t)
	}

	return visits, nil
}

// GetVisitsByIP get filtered visits by ip.
func (db *SQLite) GetVisitsByIP(ip, day string, gt, lt time.Time) (kcp.VisitsByIP, error) {
	var params []interface{}
	q := fmt.Sprint(`
	SELECT visited_at
	FROM visits
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

	rows, err := db.Query(q, params...)
	if err != nil {
		return nil, err
	}

	var t time.Time
	visits := make(kcp.VisitsByIP)
	for rows.Next() {
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		visits[ip] = append(visits[ip], t)
	}
	return visits, nil
}

// SQLiteConn returns connection to SQLite db or an error
func SQLiteConn() (*sql.DB, error) {
	os.Remove("./kcp.db")

	db, err := sql.Open("sqlite3", "./kcp.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := initSQLite(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initSQLite(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE visits (
		id integer not null primary key,
		ip text,
		day text,
		visited_at TIMESTAMP
		);`
	if _, err := db.Exec(sqlStmt); err != nil {
		fmt.Println(sqlStmt)
		fmt.Println(err)
		return err
	}
	return nil
}
