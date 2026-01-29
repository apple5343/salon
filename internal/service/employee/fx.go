package employee

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"service",
		fx.Provide(
			NewService,
		),
	)
}
