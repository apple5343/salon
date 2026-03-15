package redis

import (
	"context"
	"encoding/json"
	"errors"
	"salon/internal/repository/models"

	service "salon/internal/models"
	repo "salon/internal/repository/errors"

	"github.com/redis/go-redis/v9"
)

func (c *supplierCache) GetByID(ctx context.Context, id string) (*service.Supplier, error) {
	key := "supplier:" + id
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}
	var supplier models.Supplier
	if err = json.Unmarshal([]byte(val), &supplier); err != nil {
		return nil, err
	}
	return models.SupplierToService(&supplier), nil
}
