package revoked

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrExists          = errors.New("key already exists")
	ErrInvalidKey      = errors.New("invalid key")
	ErrInvalidDuration = errors.New("invalid expiration date")
)

type List interface {
	Create(ctx context.Context, jti string, exp time.Duration) error
	Find(ctx context.Context, jti string) (bool, error)
}
