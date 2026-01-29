package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Postgres struct {
	DSN string `env:"POSTGRES_DSN" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

func PostgresConfig() (*Postgres, error) {
	cfg := &Postgres{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
