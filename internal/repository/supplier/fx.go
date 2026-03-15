package supplier

import (
	"salon/internal/repository/supplier/postgres"
	"salon/internal/repository/supplier/redis"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"supplier",
		fx.Provide(
			postgres.NewRepository,
			redis.NewRepository,
		),
	)
}
