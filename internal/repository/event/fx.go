package event

import (
	"salon/internal/repository/event/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"event",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
