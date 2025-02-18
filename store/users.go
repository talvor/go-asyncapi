package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	ID                   uuid.UUID `db:"id"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hashed_password"`
	CreatedAt            time.Time `db:"created_at"`
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	const dml = `INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *`

	var user User
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing password: %w", err)
	}

	if err := s.db.GetContext(ctx, &user, dml, email, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return &user, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	const query = `SELECT * FROM users WHERE email = $1`

	var user User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("failed to get user by email: %s %w", email, err)
	}
	return &user, nil
}

func (s *UserStore) ByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	var user User
	if err := s.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %s %w", userID, err)
	}
	return &user, nil
}

func (u *User) ComparePassword(password string) error {
	err := CheckPasswordHash(password, u.HashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)
	return hashedPasswordBase64, nil
}

func CheckPasswordHash(password, hash string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err
}
