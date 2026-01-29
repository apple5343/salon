package http

import (
	"context"
	"log/slog"
	"salon/internal/config"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"http",
		fx.Provide(
			config.HttpServerConfig,
			NewServer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, s *Server) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := s.Start(ctx); err != nil {
								slog.Error(err.Error())
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return s.Stop(ctx)
					},
				})
			},
		),
	)
}
