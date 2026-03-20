package redis

import (
	"context"
	"encoding/json"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"time"
)

func (c *brandCache) SetByID(ctx context.Context, b *service.Brand, ttl time.Duration) error {
	key := "brand:" + b.ID
	data, err := json.Marshal(models.BrandToDatabase(b))
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}
