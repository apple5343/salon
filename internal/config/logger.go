package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Logger struct {
	Level string `env:"LOGGER_LEVEL" env-default:"dev"`
}

func LoggerConfig() (*Logger, error) {
	cfg := &Logger{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
