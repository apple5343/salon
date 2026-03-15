package app

import (
	"salon/internal/infrastructure/postgres"
	"salon/internal/infrastructure/redis"

	"go.uber.org/fx"
)

func InfrastructureModules() fx.Option {
	return fx.Module(
		"infrastructure",
		postgres.NewModule(),
		redis.NewModule(),
	)
}
