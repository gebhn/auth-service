package store

import (
	"context"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type TokenStore struct {
	conn    sqlc.DBTX
	queries *sqlc.Queries
}

func NewTokenStore(conn sqlc.DBTX) *TokenStore {
	return &TokenStore{
		conn:    conn,
		queries: sqlc.New(conn),
	}
}

func (s *TokenStore) Create(ctx context.Context, params sqlc.CreateTokenParams) (*sqlc.Token, error) {
	return s.queries.CreateToken(ctx, params)
}

func (s *TokenStore) FindOne(ctx context.Context, params sqlc.CreateTokenParams) (*sqlc.Token, error) {
	if params.Jti == "" {
		return nil, ErrInvalidInput
	}
	return s.queries.GetTokenByJTI(ctx, params.Jti)
}

func (s *TokenStore) FindMany(ctx context.Context, params sqlc.CreateTokenParams) ([]*sqlc.Token, error) {
	if params.UserID == "" {
		return nil, ErrInvalidInput
	}
	return s.queries.GetTokensForUser(ctx, params.UserID)
}
