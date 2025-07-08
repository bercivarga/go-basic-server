package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Connect() (*sql.DB, error)
	Close() error
}
