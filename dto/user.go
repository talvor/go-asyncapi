package dto

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   uuid.UUID `db:"id"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hashed_password"`
	CreatedAt            time.Time `db:"created_at"`
}

func (u *User) ComparePassword(password string) error {
	err := CheckPasswordHash(password, u.HashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}

func CheckPasswordHash(password, hash string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)
	return hashedPasswordBase64, nil
}
