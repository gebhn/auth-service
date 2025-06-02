package store

import (
	"context"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type userStore[P any, T any] struct {
	queries   *sqlc.Queries
	createFn  func(ctx context.Context, params P) (*T, error)
	findOneFn func(ctx context.Context, params P) (*T, error)
}

func NewUserStore(conn sqlc.DBTX) *userStore[sqlc.CreateUserParams, sqlc.User] {
	s := &userStore[sqlc.CreateUserParams, sqlc.User]{
		queries: sqlc.New(conn),
	}
	s.createFn = s.create
	s.findOneFn = s.findOne

	return s
}

func (s *userStore[P, T]) Create(ctx context.Context, params P) (*T, error) {
	return s.createFn(ctx, params)
}

func (s *userStore[P, T]) FindOne(ctx context.Context, params P) (*T, error) {
	return s.findOneFn(ctx, params)
}

func (s *userStore[P, T]) FindMany(ctx context.Context, params P) ([]*T, error) {
	return nil, ErrNotImplemented
}

func (s *userStore[P, T]) create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
	return s.queries.CreateUser(ctx, params)
}

func (s *userStore[P, T]) findOne(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
	switch {
	case params.Email != "":
		return s.queries.GetUserByEmail(ctx, params.Email)
	case params.Username != "":
		return s.queries.GetUserByUsername(ctx, params.Username)
	case params.UserID != "":
		return s.queries.GetUserByID(ctx, params.UserID)
	default:
		return nil, ErrInvalidInput
	}
}
