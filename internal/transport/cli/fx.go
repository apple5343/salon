package cli

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"cli",
		fx.Provide(NewCLI),
	)
}
