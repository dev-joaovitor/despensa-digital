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
	MailSMTPHost string `env:"EMAIL_SMTP_HOST" envRequired:"true"`
	MailSMTPPort string `env:"EMAIL_SMTP_PORT" envRequired:"true"`
	MailSMTPUser string `env:"EMAIL_SMTP_USER" envRequired:"true"`
	MailSMTPPassword string `env:"EMAIL_SMTP_PASSWORD" envRequired:"true"`
	MailHTTPHost string `env:"EMAIL_HTTP_HOST" envRequired:"true"`
	MailHTTPPort string `env:"EMAIL_HTTP_PASSWORD" envRequired:"true"`
	ClientAppURL string `env:"CLIENT_APP_URL" envRequired:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parsing env configurations: %w", err)
	}

	return &cfg, nil
}
