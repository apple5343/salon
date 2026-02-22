package model

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"model",
		fx.Provide(
			NewService,
		),
	)
}
