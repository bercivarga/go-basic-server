package session

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES (?, ?, ?)
	`, userID, token, expiresAt)
	return err
}

func (s *Store) IsValid(ctx context.Context, userID int64, token string) bool {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sessions
		WHERE user_id = ? AND token = ? AND expires_at > CURRENT_TIMESTAMP
	`, userID, token).Scan(&count)
	return err == nil && count == 1
}
