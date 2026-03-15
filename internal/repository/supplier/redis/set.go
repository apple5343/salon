package redis

import (
	"context"
	"encoding/json"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"time"
)

func (c *supplierCache) SetByID(ctx context.Context, s *service.Supplier, ttl time.Duration) error {
	key := "supplier:" + s.ID
	data, err := json.Marshal(models.SupplierToDatabase(s))
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}
