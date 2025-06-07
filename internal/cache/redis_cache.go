package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	c *redis.Client
}

func NewRedisCache(addr string, password string) *redisCache {
	return &redisCache{
		c: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		}),
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value string, exp time.Duration) (string, error) {
	if key == "" || value == "" || exp <= 0 {
		return "", ErrInvalidInput
	}
	return r.c.Set(ctx, key, value, exp).Result()
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrInvalidInput
	}
	return r.c.Get(ctx, key).Result()
}
