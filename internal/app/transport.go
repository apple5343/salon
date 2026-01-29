package app

import (
	"salon/internal/transport/http"

	"go.uber.org/fx"
)

func TransportModules() fx.Option {
	return fx.Module(
		"transport",
		http.NewModule(),
	)
}
