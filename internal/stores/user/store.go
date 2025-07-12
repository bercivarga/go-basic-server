package user

import (
	"context"
	"database/sql"

	"github.com/bercivarga/go-basic-server/internal/db/sqlc"
)

type Store struct {
	q *sqlc.Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{q: sqlc.New(db)}
}

func (s *Store) GetAll(ctx context.Context, limit, offset int64) ([]sqlc.User, error) {
	rows, err := s.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]sqlc.User, len(rows))
	for i, r := range rows {
		out[i] = sqlc.User{ID: r.ID, Email: r.Email}
	}
	return out, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*sqlc.User, error) {
	r, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &sqlc.User{ID: r.ID, Email: r.Email, Role: r.Role}, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	r, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &sqlc.User{ID: r.ID, Email: r.Email, PasswordHash: r.PasswordHash}, nil
}

func (s *Store) Create(ctx context.Context, email string, passwordHash string) (*sqlc.User, error) {
	r, err := s.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}
	return &sqlc.User{ID: r.ID, Email: r.Email}, nil
}

func (s *Store) GetRole(ctx context.Context, userID int64) (string, error) {
	userRole, err := s.q.GetRole(ctx, userID)
	if err != nil {
		return "", err
	}
	return userRole, nil
}
