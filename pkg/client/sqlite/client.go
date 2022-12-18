package sqlite

import (
	"database/sql"
)

type sqliteConfig struct {
	DbPath string
}

var DB *sql.DB

// NewSqliteConfig creates new sqlite config instance
func NewSqliteConfig(dbPath string) *sqliteConfig {
	return &sqliteConfig{
		DbPath: dbPath,
	}
}

// NewClient
func NewClient(config *sqliteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
