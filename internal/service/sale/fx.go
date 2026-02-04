package sale

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"sale",
		fx.Provide(
			NewService,
		),
	)
}
