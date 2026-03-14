package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type JWT struct {
	RefreshSecret string        `env:"REFRESH_SECRET" env-required:"true"`
	RefreshTTL    time.Duration `env:"REFRESH_TTL" env-default:"24h"`
	AccessSecret  string        `env:"ACCESS_SECRET" env-required:"true"`
	AccessTTL     time.Duration `env:"ACCESS_TTL" env-default:"15m"`
}

func JWTConfig() (*JWT, error) {
	cfg := &JWT{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
