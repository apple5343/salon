package simulation

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StartDate time.Time `env:"SIMULATION_START_DATE" env-default:"2022-01-01" env-layout:"2006-01-02"`
	DaysCount int       `env:"SIMULATION_DAYS_COUNT" env-default:"365"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
