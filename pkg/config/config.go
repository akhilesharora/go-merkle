package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerHost string `env:"SERVER_HOST" env-default:"localhost" env-description:"Host for the server"`
	ServerPort int    `env:"SERVER_PORT" env-default:"8080" env-description:"Port for the server"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"info" env-description:"Logging level"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	return cfg, nil
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}
