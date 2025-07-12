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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	if err != nil {
		return TokenPair{}, errors.New("token generation failed")
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, errors.New("refresh token generation failed")
	}

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

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error) {
	session, err := s.SessionStore.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return TokenPair{}, errors.New("invalid or expired refresh token")
	}

	accessToken, err := s.JwtManager.Generate(session.UserID)
	if err != nil {
		return TokenPair{}, errors.New("token generation failed")
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, errors.New("refresh token generation failed")
	}

	if err := s.SessionStore.DeleteByRefreshToken(ctx, refreshToken); err != nil {
		return TokenPair{}, errors.New("session cleanup failed")
	}

	accessTokenExpireAt, refreshTokenExpireAt := s.JwtManager.CreateExpiry()

	err = s.SessionStore.Create(ctx, session.UserID, accessToken, newRefreshToken, accessTokenExpireAt, refreshTokenExpireAt)
	if err != nil {
		return TokenPair{}, errors.New("session creation failed")
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, accessToken string) error {
	err := s.SessionStore.DeleteByToken(ctx, accessToken)
	if err != nil {
		return errors.New("could not delete session")
	}
	return nil
}
