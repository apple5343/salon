package simulation

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StartDate time.Time `env:"SIMULATION_START_DATE" env-default:"2022-01-01" env-layout:"2006-01-02"`
	DaysCount int       `env:"SIMULATION_DAYS_COUNT" env-default:"365"`

	BrandsPercent    int `env:"SIMULATION_BRANDS_PERCENT" env-default:"50"`
	ModelsPercent    int `env:"SIMULATION_MODELS_PERCENT" env-default:"50"`
	SuppliersPercent int `env:"SIMULATION_SUPPLIERS_PERCENT" env-default:"50"`
	AdminsCount      int `env:"SIMULATION_ADMINS_COUNT" env-default:"3"`
	EmployeesCount   int `env:"SIMULATION_EMPLOYEES_COUNT" env-default:"5"`

	NewEmloyeesRatio  float64 `env:"SIMULATION_NEW_EMPLOYEES_RATIO" env-default:"0.01"`
	NewCarsRatio      float64 `env:"SIMULATION_NEW_CARS_RATIO" env-default:"0.3"`
	NewBrandsRatio    float64 `env:"SIMULATION_NEW_BRANDS_RATIO" env-default:"0.02"`
	NewModelsRatio    float64 `env:"SIMULATION_NEW_MODELS_RATIO" env-default:"0.07"`
	NewSuppliersRatio float64 `env:"SIMULATION_NEW_SUPPLIERS_RATIO" env-default:"0.04"`
	NewSalesRatio     float64 `env:"SIMULATION_NEW_SALES_RATIO" env-default:"0.5"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
