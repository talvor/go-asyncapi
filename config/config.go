package config

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type Env string

const (
	EnvTest Env = "test"
	EnvDev  Env = "dev"
)

type Config struct {
	DatabaseName     string `mapstructure:"DB_NAME"`
	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseTestPort string `mapstructure:"DB_PORT_TEST"`
	DatabaseUser     string `mapstructure:"DB_USER"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	Env              Env    `mapstructure:"ENV" default:"dev"`
	ProjectRoot      string `mapstructure:"PROJECT_ROOT"`
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

	viper.AddConfigPath("..")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	defaults.SetDefaults(&cfg)

	return &cfg, nil
}
