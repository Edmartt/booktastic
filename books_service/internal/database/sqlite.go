package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct{}

func (s SQLite) GetConnection() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "test.db")

	if err != nil {
		log.Fatal("SQLITE ERROR: ", err)
	}

	return db
}
