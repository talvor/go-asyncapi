package apiserver

import (
	"log/slog"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/talvor/asyncapi/config"
	"github.com/talvor/asyncapi/store"
)

type APIServer struct {
	config *config.Config
	store  *store.Store
}

func New(config *config.Config, store *store.Store) *APIServer {
	return &APIServer{
		config: config,
		store:  store,
	}
}

func (s *APIServer) Start() error {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format:     "${time} ${locals:requestid} ${status} - ${method} ${path}\u200b\n",
		TimeFormat: time.RFC3339,
	}))

	app.Get("/ping", s.ping())

	app.Post("/auth/signup", s.signupHandler())

	host := net.JoinHostPort(s.config.APIHost, s.config.APIPort)
	slog.Info("starting server", "host", host)
	return app.Listen(net.JoinHostPort(s.config.APIHost, s.config.APIPort))
}
