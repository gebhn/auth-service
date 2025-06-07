package cache

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidInput = errors.New("invalid input")

type Cache interface {
	Set(ctx context.Context, key string, value string, exp time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
}

var _ Cache = (*redisCache)(nil)
