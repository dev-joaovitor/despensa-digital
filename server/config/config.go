package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port string `env:"APP_PORT" envRequired:"true"`
	DatabaseURL string `env:"DATABASE_URL" envRequired:"true"`
	RedisURL    string `env:"REDIS_URL" envRequired:"true"`
	BcryptRounds int `env:"BCRYPT_ROUNDS" envRequired:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parsing env configurations: %w", err)
	}

	return &cfg, nil
}
