package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddress     string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseDSN    string `env:"DATABASE_DSN" envDefault:"postgresql://postgres:postgres@postgres:5432/praktikum?sslmode=disable"`
	MigrationsPath string `env:"MIGRATIONS_PATH" envDefault:"file://migrations"`
	JWTSecret      string `env:"JWT_SECRET" envDefault:"secret"`
}

func Parse() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
