package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type sqlStore struct {
	db *sql.DB
	*sqlc.Queries
}

func NewSqlStore(db *sql.DB) *sqlStore {
	return &sqlStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (s *sqlStore) ExecTx(ctx context.Context, fn func(store Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := s.newTxStore(tx)

	if err := fn(txStore); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

func (s *sqlStore) CreateUser(ctx context.Context, p sqlc.CreateUserParams) error {
	if p.UserID == "" || p.Username == "" || p.Email == "" || p.PasswordHash == "" {
		return ErrInvalidInput
	}
	return s.Queries.CreateUser(ctx, p)
}

func (s *sqlStore) UpdateUser(ctx context.Context, p sqlc.UpdateUserParams) error {
	invalidUsername := p.Username == "" || p.Username == nil
	invalidEmail := p.Email == "" || p.Email == nil
	invalidPass := p.PasswordHash == "" || p.PasswordHash == nil

	if p.UserID == "" || (invalidUsername && invalidEmail && invalidPass) {
		return ErrInvalidInput
	}
	return s.Queries.UpdateUser(ctx, p)
}

func (s *sqlStore) GetUserByID(ctx context.Context, userID string) (*sqlc.User, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	return s.Queries.GetUserByID(ctx, userID)
}

func (s *sqlStore) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	if email == "" {
		return nil, ErrInvalidInput
	}
	return s.Queries.GetUserByEmail(ctx, email)
}

func (s *sqlStore) GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error) {
	if username == "" {
		return nil, ErrInvalidInput
	}
	return s.Queries.GetUserByUsername(ctx, username)
}

func (s *sqlStore) CreateToken(ctx context.Context, p sqlc.CreateTokenParams) error {
	if p.Jti == "" || p.UserID == "" || p.Kind == "" || p.TokenHash == "" {
		return ErrInvalidInput
	}
	if p.IssuedAt.After(time.Now()) {
		return ErrInvalidInput
	}
	if p.ExpiresAt.Before(time.Now()) {
		return ErrInvalidInput
	}
	return s.Queries.CreateToken(ctx, p)
}

func (s *sqlStore) GetTokenByJTI(ctx context.Context, jti string) (*sqlc.Token, error) {
	if jti == "" {
		return nil, ErrInvalidInput
	}
	return s.Queries.GetTokenByJTI(ctx, jti)
}

func (s *sqlStore) GetTokensForUser(ctx context.Context, userID string) ([]*sqlc.Token, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	tokens, err := s.Queries.GetTokensForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, sql.ErrNoRows
	}
	return tokens, nil
}

func (s *sqlStore) newTxStore(tx *sql.Tx) *sqlStore {
	return &sqlStore{
		db:      s.db,
		Queries: sqlc.New(tx),
	}
}
