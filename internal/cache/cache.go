package cache

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrInvalidInput = errors.New("invalid input")

type Cache interface {
	io.Closer
	Set(ctx context.Context, key string, value string, exp time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
}

var _ Cache = (*redisCache)(nil)
