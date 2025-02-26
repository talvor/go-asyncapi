package apiserver

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/talvor/asyncapi/store"
)

func AuthMiddleware(jwtManager *JwtManager, userStore *store.UserStore) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sendUnauthorized := func(message string) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": message})
		}

		var tokenString string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			tokenString = strings.TrimPrefix(authorization, "Bearer ")
		}

		if tokenString == "" {
			return sendUnauthorized("You are not logged in")
		}

		parsedToken, err := jwtManager.Parse(tokenString)
		if err != nil {
			slog.Error("failed to parse token", "error", err)
			return sendUnauthorized("You are not logged in")
		}

		if !jwtManager.IsAccessToken(parsedToken) {
			return sendUnauthorized("Not an access token")
		}

		userID, err := jwtManager.GetUserIDFromToken(parsedToken)
		if err != nil {
			slog.Error("failed to convert subject to UUID", "error", err)
			return sendUnauthorized("You are not logged in")
		}

		user, err := userStore.ByID(c.Context(), userID)
		if err != nil {
			slog.Error("failed to get user by id", "error", err)
			return sendUnauthorized("You are not logged in")
		}

		c.Locals("user", user)

		return c.Next()
	}

}
