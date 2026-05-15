package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// PingDB test database connection
func PingDB(db *sqlx.DB) {

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	name := db.DriverName()
	log.Printf("Connected to %s", name)
}
