package app

import (
	"database/sql"
	"log/slog"
)

type App struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewApp(db *sql.DB, logger *slog.Logger) *App {
	return &App{DB: db, Logger: logger}
}
