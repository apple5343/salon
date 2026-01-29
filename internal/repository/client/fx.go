package client

import (
	"salon/internal/repository/client/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"client",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
