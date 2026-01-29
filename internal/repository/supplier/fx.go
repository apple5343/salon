package supplier

import (
	"salon/internal/repository/supplier/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"supplier",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
