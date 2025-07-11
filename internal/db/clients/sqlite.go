package clients

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SQLiteDriver = "sqlite3"
)

type SQLite struct {
	Driver string
	DSN    string
	DB     *sql.DB
}

func NewSQLite(dsn string) *SQLite {
	return &SQLite{Driver: SQLiteDriver, DSN: dsn}
}

func (s *SQLite) Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", s.DSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	s.DB = db
	return db, nil
}

func (s *SQLite) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
