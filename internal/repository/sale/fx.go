package sale

import (
	"salon/internal/repository/sale/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"sale",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
