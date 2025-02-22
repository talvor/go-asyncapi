package apiserver

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrWithStatus struct {
	status int
	err    error
}

func (e ErrWithStatus) Error() string {
	return e.err.Error()
}

func NewErrWithStatus(status int, err error) error {
	return &ErrWithStatus{
		status: status,
		err:    err,
	}
}

func handler(f func(*fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := f(c); err != nil {
			status := fiber.StatusInternalServerError
			msg := http.StatusText(status)

			if e, ok := err.(*ErrWithStatus); ok {
				status = e.status
				msg = http.StatusText(e.status)
				if status == http.StatusBadRequest || status == http.StatusConflict {
					msg = e.err.Error()
				}
			}

			slog.Error("error executing handler", "error", err, "status", status, "message", msg)

			if err := c.Status(status).JSON(APIResponse[struct{}]{
				Message: msg,
			}); err != nil {
				slog.Error("error sending response", "error", err)
			}
		}
		return nil
	}
}

func encode[T any](v T, status int, c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "application/json; charset=utf-8")
	c.Status(status)
	if err := c.JSON(v); err != nil {
		return fmt.Errorf("error encoding response: %w", err)
	}
	return nil
}

type Validator interface {
	Validate() error
}

func decode[T Validator](c *fiber.Ctx) (T, error) {
	var t T
	if err := c.BodyParser(&t); err != nil {
		return t, fmt.Errorf("decoding request body: %w", err)
	}

	if err := t.Validate(); err != nil {
		return t, err
	}

	return t, nil
}
