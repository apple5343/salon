package config

import (
	"fmt"
	"salon/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

func LoggerConfig() (*logger.Config, error) {
	cfg := &logger.Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
