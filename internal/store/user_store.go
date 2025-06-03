package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gebhn/auth-service/internal/db/sqlc"
)

type UserStore struct {
	conn    sqlc.DBTX
	queries *sqlc.Queries
}

func NewUserStore(conn sqlc.DBTX) *UserStore {
	return &UserStore{
		conn:    conn,
		queries: sqlc.New(conn),
	}
}

func (s *UserStore) Create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
	return s.queries.CreateUser(ctx, params)
}

func (s *UserStore) FindOne(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
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

func (s *UserStore) FindMany(ctx context.Context, params sqlc.CreateUserParams) ([]*sqlc.User, error) {
	return nil, ErrNotImplemented
}

func (s *UserStore) WithTransaction(ctx context.Context, fn func(ts *UserStore) error) error {
	conn, ok := s.conn.(*sql.DB)
	if !ok {
		return fmt.Errorf("cannot start a new transaction from an existing transaction store")
	}

	tx, err := conn.BeginTx(ctx, nil) // Start a new transaction
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ts := NewUserStore(tx)

	defer func() {
		if r := recover(); r != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("panic during transaction, rollback failed: %v, original panic: %v", rbErr, r)
			} else {
				log.Printf("prror during transaction, rolled back: %v", r)
			}
			panic(r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("error during transaction, rollback failed: %v, original error: %v", rbErr, err)
			} else {
				log.Printf("error during transaction, rolled back: %v", err)
			}
		}
	}()

	err = fn(ts)
	if err != nil {
		return err
	}

	return tx.Commit()
}
