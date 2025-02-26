package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/talvor/asyncapi/dto"
)

type RefreshTokenStore struct {
	db *sqlx.DB
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *RefreshTokenStore) Create(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*dto.RefreshToken, error) {
	const dml = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3) RETURNING *`

	hashedToken, err := dto.HashToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing token: %w", err)
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expires at from token")
	}

	var refreshToken dto.RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, dml, userID, hashedToken, expiresAt.Time); err != nil {
		return nil, fmt.Errorf("failed to insert refresh token record: %w", err)
	}
	return &refreshToken, nil
}

func (s *RefreshTokenStore) ByPrimaryKey(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*dto.RefreshToken, error) {
	const query = `SELECT * FROM refresh_tokens WHERE user_id = $1 AND hashed_token = $2;`

	hashedToken, err := dto.HashToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to hashing token: %w", err)
	}

	var refreshToken dto.RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, query, userID, hashedToken); err != nil {
		return nil, fmt.Errorf("failed to get refresh token record: %w", err)
	}
	return &refreshToken, nil
}

func (s *RefreshTokenStore) DeleteUserTokens(ctx context.Context, userID uuid.UUID) (sql.Result, error) {
	const dml = `DELETE FROM refresh_tokens WHERE user_id = $1;`

	result, err := s.db.ExecContext(ctx, dml, userID)
	if err != nil {
		return result, fmt.Errorf("failed to delete user refresh tokens record: %w", err)
	}
	return result, nil
}

func (s *RefreshTokenStore) DeleteUserTokensThenCreate(ctx context.Context, userID uuid.UUID, token *jwt.Token) (*dto.RefreshToken, error) {
	_, err := s.DeleteUserTokens(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}

	return s.Create(ctx, userID, token)
}
