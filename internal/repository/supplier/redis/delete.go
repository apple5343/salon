package redis

import (
	"context"
)

func (c *supplierCache) DeleteByID(ctx context.Context, id string) error {
	key := "supplier:" + id
	return c.client.Del(ctx, key).Err()
}
