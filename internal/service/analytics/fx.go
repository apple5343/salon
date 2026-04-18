package analytics

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"analytics",
		fx.Provide(
			NewService,
		),
	)
}
