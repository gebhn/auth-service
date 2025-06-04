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

func (l *cacheRevokedList) Create(ctx context.Context, jti string, kind pb.TokenKind exp time.Duration) error {
	if jti == "" {
		return ErrInvalidKey
	}
	if exp.Abs() <= config.GetAccessTokenDuration() {
		return ErrInvalidDuration
	}
}

func (l *cacheRevokedList) Find(ctx context.Context, jti string) (bool, error) {
	v, err := l.c.Get(ctx, jti)
	if err != nil {
		return false, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return false, err
	}
	return i > 0, nil
}
