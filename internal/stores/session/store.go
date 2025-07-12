package session

import (
	"context"
	"database/sql"
	"time"

	"github.com/bercivarga/go-basic-server/internal/db/sqlc"
)

type Store struct {
	q *sqlc.Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{q: sqlc.New(db)}
}

func (s *Store) Create(ctx context.Context, userID int64, token, refreshToken string, expiresAt, refreshExpiresAt time.Time) error {
	err := s.q.CreateSession(ctx, sqlc.CreateSessionParams{
		UserID:           userID,
		Token:            token,
		ExpiresAt:        expiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) IsValid(ctx context.Context, userID int64, token string) bool {
	count, err := s.q.IsValidSession(ctx, sqlc.IsValidSessionParams{
		UserID: userID,
		Token:  token,
	})
	if err != nil {
		return false
	}
	return count == 1
}

func (s *Store) DeleteByToken(ctx context.Context, token string) error {
	err := s.q.DeleteSessionByToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetByRefreshToken(ctx context.Context, token string) (sqlc.Session, error) {
	return s.q.GetSessionByRefreshToken(ctx, token)
}

func (s *Store) DeleteByRefreshToken(ctx context.Context, token string) error {
	return s.q.DeleteSessionByRefreshToken(ctx, token)
}
