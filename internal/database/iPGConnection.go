package database

import "github.com/jmoiron/sqlx"

type IDBConnection interface {
	GetConnection() *sqlx.DB
}
