package app

import (
	"database/sql"
	"log/slog"

	"github.com/bercivarga/go-basic-server/internal/auth"
	"github.com/bercivarga/go-basic-server/internal/config"
	"github.com/bercivarga/go-basic-server/internal/logger"
	"github.com/bercivarga/go-basic-server/internal/stores/session"
)

type App struct {
	DB           *sql.DB
	Logger       *slog.Logger
	Config       *config.Config
	SessionStore *session.Store
	JwtManager   *auth.JWTManager
}

func NewApp(db *sql.DB) *App {
	logger := logger.New()
	config := config.Load()
	sessionStore := session.NewStore(db)
	jwtManager := auth.NewJWTManager(config.JWTSecret)

	return &App{DB: db, Logger: logger, Config: config, SessionStore: sessionStore, JwtManager: jwtManager}
}
