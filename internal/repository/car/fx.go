package car

import (
	"salon/internal/repository/car/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"car",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
