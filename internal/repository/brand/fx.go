package brand

import (
	"salon/internal/repository/brand/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"brand",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
