package app

import (
	"salon/internal/config"
	"salon/pkg/logger"

	"go.uber.org/fx"
)

func Logger() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(
			config.LoggerConfig,
			logger.NewLogger,
		),
	)
}
