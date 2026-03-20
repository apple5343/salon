package redis

import (
	"salon/internal/repository"

	"github.com/redis/go-redis/v9"
)

type brandCache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) repository.BrandCache {
	return &brandCache{
		client: client,
	}
}
