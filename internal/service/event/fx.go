package event

import (
	"context"
	"salon/internal/service"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"event",
		fx.Provide(
			NewService,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, s service.EventService) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return s.Shutdown(ctx)
					},
				})
			},
		),
	)
}
