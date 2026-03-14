package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Redis struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"redis"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

func RedisConfig() (*Redis, error) {
	cfg := &Redis{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
