package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sqlx uses it indirectly
)

type Postgres struct {
	config *PGConfig
	conn   *sqlx.DB
}

func NewPostgres(cfg *PGConfig) *Postgres {
	return &Postgres{
		config: cfg,
	}
}

// GetConnection connects to specific database
func (p *Postgres) GetConnection() *sqlx.DB {

	if p.conn != nil {
		return p.conn
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", p.config.User, p.config.DB, p.config.Password, p.config.Host, p.config.Port))

	if err != nil {
		log.Fatal("CONNECTION ERROR: ", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(25)

	p.conn = db

	return p.conn
}
