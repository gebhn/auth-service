package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, exp time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
}

var _ Cache = (*redisCache)(nil)
