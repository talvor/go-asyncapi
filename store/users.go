package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/talvor/asyncapi/dto"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*dto.User, error) {
	const dml = `INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *`

	var user dto.User
	hashedPassword, err := dto.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing password: %w", err)
	}

	if err := s.db.GetContext(ctx, &user, dml, email, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return &user, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*dto.User, error) {
	const query = `SELECT * FROM users WHERE email = $1`

	var user dto.User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("failed to get user by email: %s %w", email, err)
	}
	return &user, nil
}

func (s *UserStore) ByID(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	var user dto.User
	if err := s.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %s %w", userID, err)
	}
	return &user, nil
}
