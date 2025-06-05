package store

import (
	"context"
	"errors"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

var ErrInvalidInput = errors.New("invalid input")

type Store interface {
	sqlc.Querier
	ExecTx(ctx context.Context, fn func(s Store) error) error
}

var _ Store = (*sqlStore)(nil)
