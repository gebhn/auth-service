package revoked

import (
	"context"
	"strconv"
	"time"

	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/cache"
	"github.com/gebhn/auth-service/internal/config"
)

type cacheRevokedList struct {
	c cache.Cache
}

func NewCacheRevokedList(c cache.Cache) *cacheRevokedList {
	return &cacheRevokedList{c: c}
}

func (r *cacheRevokedList) Create(ctx context.Context, jti string, kind pb.TokenKind, exp time.Duration) error {
	if jti == "" {
		return ErrInvalidKey
	}
	if exp.Abs() < config.GetTokenDuration(kind) {
		return ErrInvalidDuration
	}
	if _, err := r.c.Set(ctx, jti, "1", exp); err != nil {
		return err
	}
	return nil
}

func (r *cacheRevokedList) Find(ctx context.Context, jti string) (bool, error) {
	v, err := r.c.Get(ctx, jti)
	if err != nil {
		return false, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return false, err
	}
	return i > 0, nil
}
