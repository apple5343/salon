package redis

import (
	"context"
	"errors"

	repo "salon/internal/repository/errors"

	"github.com/redis/go-redis/v9"
)

func (c *brandCache) DeleteByID(ctx context.Context, id string) error {
	key := "supplier:" + id
	if err := c.client.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return repo.ErrNotFound
		}
		return err
	}
	return nil
}
