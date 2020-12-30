package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Gorm contains connection to db.
type Gorm struct {
	*gorm.DB
}

// Visit contains fields for visit row in db.
type Visit struct {
	VisitedAt time.Time
	IP        string
	Day       string
}

// SQLiteGormConn returns connection to gorm SQLite db or an error.
func SQLiteGormConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./kcp.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := initSQLiteGorm(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initSQLiteGorm(db *gorm.DB) error {
	return db.AutoMigrate(&Visit{})
}

// InsertEvent inserts kcp.Event into db.
func (db *Gorm) InsertEvent(e kcp.Event) error {
	return db.Create(&Visit{
		VisitedAt: e.VisitedAt,
		IP:        e.IP,
		Day:       e.Day,
	}).Error
}

// GetVisits get visits grouped by ip.
func (db *Gorm) GetVisits() (kcp.VisitsByIP, error) {
	visits := make(kcp.VisitsByIP)
	rows, err := db.Model(&Visit{}).Select("ip", "visited_at").Rows()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var visit Visit
		if err := db.ScanRows(rows, &visit); err != nil {
			fmt.Println(err)
			return nil, err
		}
		visits[visit.IP] = append(visits[visit.IP], visit.VisitedAt)
	}

	return visits, nil
}

// GetVisitsByIP get filtered visits by ip.
func (db *Gorm) GetVisitsByIP(ip, day string, gt, lt time.Time) (kcp.VisitsByIP, error) {
	var visits []Visit

	tx := db.DB
	if !gt.IsZero() {
		tx = tx.Where("visited_at > ?", gt)
	}
	if !lt.IsZero() {
		tx = tx.Where("visited_at < ?", lt)
	}
	if day != "" {
		tx = tx.Where("day = ?", day)
	}
	res := tx.Find(&visits)

	visitsByIP := make(kcp.VisitsByIP)
	for _, v := range visits {
		visitsByIP[v.IP] = append(visitsByIP[v.IP], v.VisitedAt)
	}
	return visitsByIP, res.Error
}
