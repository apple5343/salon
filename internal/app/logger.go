package app

import (
	"salon/pkg/logger"

	"go.uber.org/fx"
)

func Logger() fx.Option {
	return fx.Module(
		"logger",
		logger.NewModule(),
	)
}
