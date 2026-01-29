package supplier

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"supplier",
		fx.Provide(
			NewService,
		),
	)
}
