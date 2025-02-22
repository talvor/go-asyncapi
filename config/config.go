package config

import (
	"fmt"
	"log"

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
	DatabaseUser     string `mapstructure:"DB_USER"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	Env              Env    `mapstructure:"ENV" default:"dev"`
	ProjectRoot      string `mapstructure:"PROJECT_ROOT"`
	APIPort          string `mapstructure:"API_PORT"`
	APIHost          string `mapstructure:"API_HOST"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

func (c *Config) SetDatabasePort(port string) {
	c.DatabasePort = port
}

var config Config

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read config: %w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	defaults.SetDefaults(&config)
}

func GetConfig() *Config {
	return &config
}
