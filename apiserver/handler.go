package apiserver

import (
	"database/sql"
	"errors"
	"fmt"

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
