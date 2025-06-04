package revoked

import (
	"context"
	"errors"
	"time"

	"github.com/gebhn/auth-service/api/pb"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidKey      = errors.New("invalid key")
	ErrInvalidDuration = errors.New("invalid expiration date")
)

type List interface {
	Create(ctx context.Context, jti string, kind pb.TokenKind, exp time.Duration) error
	Find(ctx context.Context, jti string) (bool, error)
}

var _ List = (*cacheRevokedList)(nil)
