package store

import (
	"context"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type revokedTokenStore[P any, T any] struct {
	queries    *sqlc.Queries
	createFn   func(ctx context.Context, params P) (*T, error)
	findOneFn  func(ctx context.Context, params P) (*T, error)
	findManyFn func(ctx context.Context, params P) ([]*T, error)
}

func NewRevokedTokenStore(conn sqlc.DBTX) *revokedTokenStore[sqlc.CreateTokenParams, string] {
	s := &revokedTokenStore[sqlc.CreateTokenParams, string]{
		queries: sqlc.New(conn),
	}
	s.createFn = s.create
	s.findOneFn = s.findOne
	s.findManyFn = s.findMany

	return s
}

func (s *revokedTokenStore[P, T]) Create(ctx context.Context, params P) (*T, error) {
	return s.createFn(ctx, params)
}

func (s *revokedTokenStore[P, T]) FindOne(ctx context.Context, params P) (*T, error) {
	return s.findOneFn(ctx, params)
}

func (s *revokedTokenStore[P, T]) FindMany(ctx context.Context, params P) ([]*T, error) {
	return s.findManyFn(ctx, params)
}

func (s *revokedTokenStore[P, T]) create(ctx context.Context, params sqlc.CreateTokenParams) (*string, error) {
	res, err := s.queries.CreateRevokedToken(ctx, params.Jti)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *revokedTokenStore[P, T]) findOne(ctx context.Context, params sqlc.CreateTokenParams) (*string, error) {
	if params.Jti == "" {
		return nil, ErrInvalidInput
	}
	res, err := s.queries.GetRevokedTokenByJti(ctx, params.Jti)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *revokedTokenStore[P, T]) findMany(ctx context.Context, params sqlc.CreateTokenParams) ([]*string, error) {
	if params.UserID == "" {
		return nil, ErrInvalidInput
	}
	revoked, err := s.queries.GetRevocableTokensByUser(ctx, params.UserID)
	if err != nil {
		return nil, err
	}
	result := make([]*string, 0, len(revoked))
	for _, jti := range revoked {
		result = append(result, &jti)
	}
	return result, nil
}
