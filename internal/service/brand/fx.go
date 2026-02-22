package brand

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"brand",
		fx.Provide(
			NewService,
		),
	)
}
