package store

import (
	"context"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type revokedTokenStore struct {
	conn    sqlc.DBTX
	queries *sqlc.Queries
}

func NewRevokedTokenStore(conn sqlc.DBTX) *revokedTokenStore {
	return &revokedTokenStore{
		conn:    conn,
		queries: sqlc.New(conn),
	}
}

func (s *revokedTokenStore) Create(ctx context.Context, params sqlc.CreateTokenParams) (string, error) {
	res, err := s.queries.CreateRevokedToken(ctx, params.Jti)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *revokedTokenStore) FindOne(ctx context.Context, params sqlc.CreateTokenParams) (string, error) {
	if params.Jti == "" {
		return "", ErrInvalidInput
	}
	res, err := s.queries.GetRevokedTokenByJti(ctx, params.Jti)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *revokedTokenStore) FindMany(ctx context.Context, params sqlc.CreateTokenParams) ([]string, error) {
	if params.UserID == "" {
		return nil, ErrInvalidInput
	}
	revoked, err := s.queries.GetRevocableTokensByUser(ctx, params.UserID)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(revoked))
	result = append(result, revoked...)
	return result, nil
}
