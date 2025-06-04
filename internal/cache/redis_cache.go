package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gebhn/auth-service/internal/config"
)

type redisCache struct {
	c *redis.Client
}

func NewRedisCache() *redisCache {
	return &redisCache{
		c: redis.NewClient(&redis.Options{
			Addr:     config.GetRedisAddress(),
			Password: config.GetRedisPassword(),
		}),
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value string, exp time.Duration) (string, error) {
	return r.c.Set(ctx, key, value, exp).Result()
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.c.Get(ctx, key).Result()
}
