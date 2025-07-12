package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bercivarga/go-basic-server/internal/auth"
	"github.com/bercivarga/go-basic-server/internal/config"
	"github.com/bercivarga/go-basic-server/internal/stores/session"
	"github.com/bercivarga/go-basic-server/internal/stores/user"
	"github.com/bercivarga/go-basic-server/internal/utils"
)

type Service struct {
	UserStore    *user.Store
	SessionStore *session.Store
	JwtManager   *auth.JWTManager
}

func New(db *sql.DB, config *config.Config) *Service {
	userStore := user.NewStore(db)
	sessionStore := session.NewStore(db)
	jwtManager := auth.NewJWTManager(config.JWTSecret)
	return &Service{UserStore: userStore, SessionStore: sessionStore, JwtManager: jwtManager}
}

func (s *Service) CheckRole(ctx context.Context, userID int64, expected string) error {
	role, err := s.UserStore.GetRole(ctx, userID)
	if err != nil {
		return err
	}
	if role != expected {
		return errors.New("forbidden")
	}
	return nil
}

func (s *Service) GetRole(ctx context.Context, userID int64) (string, error) {
	return s.UserStore.GetRole(ctx, userID)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *Service) Login(ctx context.Context, email, password string) (TokenPair, error) {
	user, err := s.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return TokenPair{}, errors.New("user not found")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return TokenPair{}, errors.New("invalid password")
	}

	accessToken, err := s.JwtManager.Generate(user.ID)
	refreshToken, err := utils.GenerateRefreshToken()
	accessExp, refreshExp := s.JwtManager.CreateExpiry()

	err = s.SessionStore.Create(ctx, user.ID, accessToken, refreshToken, accessExp, refreshExp)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
