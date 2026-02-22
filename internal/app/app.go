package app

import (
	"salon/pkg/clock"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		Config(),
		Logger(),
		InfrastructureModules(),
		RepositoryModules(),
		ServiceModules(),
		TransportModules(),
		clock.ClockModule(),
	)
}
