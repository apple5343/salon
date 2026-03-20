package brand

import (
	"salon/internal/repository/brand/postgres"
	"salon/internal/repository/brand/redis"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"brand",
		fx.Provide(
			postgres.NewRepository,
			redis.NewCache,
		),
	)
}
