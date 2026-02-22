package simulation

import "go.uber.org/fx"

func SimulationModule() fx.Option {
	return fx.Module(
		"simulation",
		fx.Provide(
			NewConfig,
			NewSimulation,
		),
	)
}
