package store

import (
	"context"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type tokenStore[P any, T any] struct {
	queries    *sqlc.Queries
	createFn   func(ctx context.Context, params P) (*T, error)
	findOneFn  func(ctx context.Context, params P) (*T, error)
	findManyFn func(ctx context.Context, params P) ([]*T, error)
}

func NewTokenStore(conn sqlc.DBTX) *tokenStore[sqlc.CreateTokenParams, sqlc.Token] {
	s := &tokenStore[sqlc.CreateTokenParams, sqlc.Token]{
		queries: sqlc.New(conn),
	}
	s.createFn = s.create
	s.findOneFn = s.findOne
	s.findManyFn = s.findMany

	return s
}

func (s *tokenStore[P, T]) Create(ctx context.Context, params P) (*T, error) {
	return s.createFn(ctx, params)
}

func (s *tokenStore[P, T]) FindOne(ctx context.Context, params P) (*T, error) {
	return s.findOneFn(ctx, params)
}

func (s *tokenStore[P, T]) FindMany(ctx context.Context, params P) ([]*T, error) {
	return s.findManyFn(ctx, params)
}

func (s *tokenStore[P, T]) create(ctx context.Context, params sqlc.CreateTokenParams) (*sqlc.Token, error) {
	return s.queries.CreateToken(ctx, params)
}

func (s *tokenStore[P, T]) findOne(ctx context.Context, params sqlc.CreateTokenParams) (*sqlc.Token, error) {
	if params.Jti == "" {
		return nil, ErrInvalidInput
	}
	return s.queries.GetTokenByJTI(ctx, params.Jti)
}

func (s *tokenStore[P, T]) findMany(ctx context.Context, params sqlc.CreateTokenParams) ([]*sqlc.Token, error) {
	if params.UserID == "" {
		return nil, ErrInvalidInput
	}
	return s.queries.GetTokensForUser(ctx, params.UserID)
}
