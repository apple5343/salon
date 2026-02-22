package model

import (
	"salon/internal/repository/model/postgres"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"model",
		fx.Provide(
			postgres.NewRepository,
		),
	)
}
