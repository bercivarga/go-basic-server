package user

import (
	"database/sql"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type store struct {
	db *sql.DB
}

type UserStore interface {
	GetAllUsers() ([]User, error)
	GetUser(id string) (*User, error)
}

func NewUserStore(db *sql.DB) UserStore {
	return &store{db: db}
}

func (s *store) GetAllUsers() ([]User, error) {
	return nil, nil // TODO
}

func (s *store) GetUser(id string) (*User, error) {
	return nil, nil // TODO
}
