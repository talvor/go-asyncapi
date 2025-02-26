package dto

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID      uuid.UUID `db:"user_id"`
	HashedToken string    `db:"hashed_token"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

func HashToken(token *jwt.Token) (string, error) {
	// bytes, err := bcrypt.GenerateFromPassword([]byte(token.Raw), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", err
	// }

	h := sha256.New()
	h.Write([]byte(token.Raw))
	hashedBytes := h.Sum(nil)

	hashedTokenBase64 := base64.StdEncoding.EncodeToString(hashedBytes)
	return hashedTokenBase64, nil
}
