package store

import (
	"context"
	"errors"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotImplemented = errors.New("not implemented")
)

type Store[P any, T any] interface {
	Create(context.Context, P) (*T, error)
	FindOne(context.Context, P) (*T, error)
	FindMany(context.Context, P) ([]*T, error)
}

var (
	_ Store[sqlc.CreateUserParams, sqlc.User]   = (*userStore[sqlc.CreateUserParams, sqlc.User])(nil)
	_ Store[sqlc.CreateTokenParams, sqlc.Token] = (*tokenStore[sqlc.CreateTokenParams, sqlc.Token])(nil)
	_ Store[sqlc.CreateTokenParams, string]     = (*revokedTokenStore[sqlc.CreateTokenParams, string])(nil)
)
