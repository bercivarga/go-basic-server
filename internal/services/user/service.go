package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bercivarga/go-basic-server/internal/stores/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store *user.Store
}

func New(db *sql.DB) *Service {
	store := user.NewStore(db)
	return &Service{store: store}
}

type CreateUserRequest struct {
	Email    string
	Password string
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) error {
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("password hashing failed")
	}

	// Create the user
	_, err = s.store.Create(ctx, req.Email, string(hash))
	if err != nil {
		return errors.New("user already exists or database error")
	}

	return nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int64) (*UserResponse, error) {
	user, err := s.store.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

type ListUsersRequest struct {
	Limit  int64
	Offset int64
}

func (s *Service) ListUsers(ctx context.Context, req ListUsersRequest) ([]UserResponse, error) {
	users, err := s.store.GetAll(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}
	}

	return response, nil
}
