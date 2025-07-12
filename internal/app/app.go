package app

import (
	"database/sql"
	"log/slog"

	"github.com/bercivarga/go-basic-server/internal/config"
	"github.com/bercivarga/go-basic-server/internal/logger"
	"github.com/bercivarga/go-basic-server/internal/services/auth"
)

type App struct {
	DB          *sql.DB
	Logger      *slog.Logger
	Config      *config.Config
	AuthService *auth.Service
}

func NewApp(db *sql.DB) *App {
	logger := logger.New()
	config := config.Load()
	authService := auth.New(db, config)

	return &App{DB: db, Logger: logger, Config: config, AuthService: authService}
}
