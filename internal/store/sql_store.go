package store

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *sqlStore) ExecTx(ctx context.Context, fn func(queries *sqlc.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	q := sqlc.New(tx)

	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (s *sqlStore) CreateUser(ctx context.Context, p sqlc.CreateUserParams) error {
	if p.UserID == "" || p.Username == "" || p.Email == "" || p.PasswordHash == "" {
		return fmt.Errorf("failed to create user: %w", ErrInvalidInput)
	}
	if err := s.Queries.CreateUser(ctx, p); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *sqlStore) UpdateUser(ctx context.Context, p sqlc.UpdateUserParams) error {
	if p.UserID == "" || (p.Username == "" && p.Email == "") {
		return fmt.Errorf("failed to update user: %w", ErrInvalidInput)
	}
	if err := s.Queries.UpdateUser(ctx, p); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *sqlStore) GetUserByID(ctx context.Context, userID string) (*sqlc.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("failed to get user by id: %w", ErrInvalidInput)
	}
	user, err := s.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *sqlStore) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	if email == "" {
		return nil, fmt.Errorf("failed to get user by email: %w", ErrInvalidInput)
	}
	user, err := s.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *sqlStore) GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error) {
	if username == "" {
		return nil, fmt.Errorf("failed to get user by username: %w", ErrInvalidInput)
	}
	user, err := s.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

func (s *sqlStore) CreateToken(ctx context.Context, p sqlc.CreateTokenParams) error {
	if p.Jti == "" || p.UserID == "" || p.Kind == "" || p.TokenHash == "" {
		return fmt.Errorf("failed to create token: %w", ErrInvalidInput)
	}
	if p.IssuedAt.Compare(time.Now()) > 0 {
		return fmt.Errorf("failed to create token: %w", ErrInvalidInput)
	}
	if err := s.Queries.CreateToken(ctx, p); err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}
	return nil
}

func (s *sqlStore) GetTokenByJTI(ctx context.Context, jti string) (*sqlc.Token, error) {
	if jti == "" {
		return nil, fmt.Errorf("failed to get token by jti: %w", ErrInvalidInput)
	}
	token, err := s.Queries.GetTokenByJTI(ctx, jti)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by jti: %w", err)
	}
	return token, nil
}

func (s *sqlStore) GetTokensForUser(ctx context.Context, userID string) ([]*sqlc.Token, error) {
	if userID == "" {
		return nil, fmt.Errorf("failed to get tokens for user: %w", ErrInvalidInput)
	}
	tokens, err := s.Queries.GetTokensForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens for user %s: %w", userID, err)
	}
	return tokens, nil
}
