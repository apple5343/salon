package redis

import (
	"context"
	"salon/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"redis",
		fx.Provide(
			config.RedisConfig,
			NewClient,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, client *redis.Client) {
				lc.Append(fx.Hook{
					OnStop: func(_ context.Context) error {
						return client.Close()
					},
				})
			},
		),
	)
}
