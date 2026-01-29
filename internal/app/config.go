package app

import (
	"salon/internal/config"

	"go.uber.org/fx"
)

func Config() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(
			config.JWTConfig,
		),
	)
}
