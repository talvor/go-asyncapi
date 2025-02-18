package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvTest Env = "test"
	EnvDev  Env = "dev"
)

type Config struct {
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseTestPort string `env:"DB_PORT_TEST"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	Env              Env    `env:"ENV" envDefault:"dev"`
	ProjectRoot      string `env:"PROJECT_ROOT"`
}

func (c Config) DatabaseURL() string {
	port := c.DatabasePort
	if c.Env == EnvTest {
		port = c.DatabaseTestPort
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		port,
		c.DatabaseName,
	)
}

func (c *Config) SetDatabaseTestPort(port string) {
	c.DatabaseTestPort = port
}

func New() (*Config, error) {
	var cfg Config
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}
