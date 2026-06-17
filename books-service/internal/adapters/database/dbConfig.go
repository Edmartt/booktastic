package database

import (
	"fmt"
	"os"
)

type PGConfig struct {
	User     string
	DB       string
	Password string
	Host     string
	Port     string
}

func LoadPGConfig() (*PGConfig, error) {
	user := os.Getenv("PG_USER")

	if user == "" {
		return nil, fmt.Errorf("PG_USER environment variable required")
	}
	db := os.Getenv("PG_DB")

	if db == "" {
		return nil, fmt.Errorf("PG_DB environment variable required")
	}

	password := os.Getenv("PG_PASSWORD")

	if password == "" {
		return nil, fmt.Errorf("PG_PASSWORD environment variable required")
	}

	host := os.Getenv("HOST")

	if host == "" {
		return nil, fmt.Errorf("HOST environment variable required")
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, fmt.Errorf("PORT environment variable required")
	}

	return &PGConfig{
		User:     user,
		DB:       db,
		Password: password,
		Host:     host,
		Port:     port,
	}, nil
}
