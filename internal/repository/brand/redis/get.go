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

func (c *brandCache) GetByID(ctx context.Context, id string) (*service.Brand, error) {
	key := "brand:" + id
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}
	var brand models.Brand
	if err = json.Unmarshal([]byte(val), &brand); err != nil {
		return nil, err
	}
	b := models.BrandToService(&brand)
	return b, nil
}
