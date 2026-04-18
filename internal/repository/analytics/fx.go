package analytics

import (
	"salon/internal/repository/analytics/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"analytics",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
