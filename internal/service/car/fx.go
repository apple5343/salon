package car

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"car",
		fx.Provide(
			NewService,
		),
	)
}
