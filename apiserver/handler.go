package apiserver

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type APIResponse[T any] struct {
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignupRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *APIServer) ping() fiber.Handler {
	return handler(func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
}

func (s *APIServer) signupHandler() fiber.Handler {
	return handler(func(c *fiber.Ctx) error {
		req, err := decode[SignupRequest](c)
		if err != nil {
			return NewErrWithStatus(fiber.StatusBadRequest, err)
		}

		// Check if user already exists
		existingUser, err := s.store.Users.ByEmail(c.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}
		if existingUser != nil {
			return NewErrWithStatus(fiber.StatusConflict, fmt.Errorf("email already registered"))
		}

		_, err = s.store.Users.CreateUser(c.Context(), req.Email, req.Password)
		if err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		if err = encode(APIResponse[struct{}]{Message: "successfully signed up user"}, fiber.StatusCreated, c); err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		return nil
	})
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r SigninRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *APIServer) signinHandler() fiber.Handler {
	return handler(func(c *fiber.Ctx) error {
		req, err := decode[SigninRequest](c)
		if err != nil {
			return NewErrWithStatus(fiber.StatusBadRequest, err)
		}

		user, err := s.store.Users.ByEmail(c.Context(), req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return NewErrWithStatus(fiber.StatusNotFound, err)
		}

		if err = user.ComparePassword(req.Password); err != nil {
			return NewErrWithStatus(fiber.StatusUnauthorized, err)
		}

		tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID)
		if err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokens.DeleteUserTokensThenCreate(c.Context(), user.ID, tokenPair.RefreshToken)
		if err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		if err := encode(APIResponse[SigninResponse]{
			Data: &SigninResponse{
				AccessToken:  tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.RefreshToken.Raw,
			},
		}, fiber.StatusOK, c); err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		return nil
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *APIServer) refreshTokenHandler() fiber.Handler {
	return handler(func(c *fiber.Ctx) error {

		req, err := decode[RefreshTokenRequest](c)
		if err != nil {
			return NewErrWithStatus(fiber.StatusBadRequest, err)
		}

		currentRefreshToken, err := s.jwtManager.Parse(req.RefreshToken)
		if err != nil {
			return NewErrWithStatus(fiber.StatusUnauthorized, err)
		}

		userID, err := s.jwtManager.GetUserIDFromToken(currentRefreshToken)
		if err != nil {
			return NewErrWithStatus(fiber.StatusUnauthorized, err)
		}

		currentRefreshTokenRecord, err := s.store.RefreshTokens.ByPrimaryKey(c.Context(), userID, currentRefreshToken)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, sql.ErrNoRows) {
				status = fiber.StatusUnauthorized
			}
			return NewErrWithStatus(status, err)
		}

		if currentRefreshTokenRecord.ExpiresAt.Before(time.Now()) {
			return NewErrWithStatus(fiber.StatusUnauthorized, errors.New("refresh token expired"))
		}

		tokenPair, err := s.jwtManager.GenerateTokenPair(userID)
		if err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		_, err = s.store.RefreshTokens.DeleteUserTokensThenCreate(c.Context(), userID, tokenPair.RefreshToken)
		if err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		if err := encode(APIResponse[RefreshTokenResponse]{
			Data: &RefreshTokenResponse{
				AccessToken:  tokenPair.AccessToken.Raw,
				RefreshToken: tokenPair.RefreshToken.Raw,
			},
		}, fiber.StatusOK, c); err != nil {
			return NewErrWithStatus(fiber.StatusInternalServerError, err)
		}

		return nil
	})
}
