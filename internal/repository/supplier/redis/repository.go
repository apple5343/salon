package redis

import (
	"salon/internal/repository"

	"github.com/redis/go-redis/v9"
)

type supplierCache struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) repository.SupplierCache {
	return &supplierCache{
		client: client,
	}
}
