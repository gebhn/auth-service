package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/config"
	"github.com/gebhn/auth-service/internal/db/sqlc"
	"github.com/gebhn/auth-service/internal/store"
)

type userAuth struct {
	us *store.UserStore
	ts *store.TokenStore
}

func NewUserAuth(us *store.UserStore, ts *store.TokenStore) *userAuth {
	return &userAuth{
		us: us,
		ts: ts,
	}
}

func (a *userAuth) Register(ctx context.Context, rr *pb.RegisterRequest) (*pb.User, error) {
	switch {
	case !rr.HasUsername():
		return nil, ErrInvalidInput
	case !rr.HasEmail():
		return nil, ErrInvalidInput
	case !rr.HasPassword():
		return nil, ErrInvalidInput
	default:
		findUserByUsername, err := a.us.FindOne(ctx, sqlc.CreateUserParams{Username: rr.GetUsername()})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		if findUserByUsername.Username != "" {
			return nil, ErrExists
		}

		findUserByEmail, err := a.us.FindOne(ctx, sqlc.CreateUserParams{Email: rr.GetEmail()})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		if findUserByEmail.Email != "" {
			return nil, ErrExists
		}

		hash, err := a.hashPassword(rr.GetPassword())
		if err != nil {
			return nil, err
		}
		in, err := a.us.Create(ctx, sqlc.CreateUserParams{
			UserID:       uuid.New().String(),
			Username:     rr.GetUsername(),
			Email:        rr.GetEmail(),
			PasswordHash: hash,
		})
		if err != nil {
			return nil, err
		}

		var out pb.User

		out.SetEmail(in.Email)
		out.SetUserId(in.UserID)
		out.SetUsername(in.Username)

		return &out, nil
	}
}

func (a *userAuth) Login(ctx context.Context, lr *pb.LoginRequest) ([2]*pb.Token, error) {
	var tokens [2]*pb.Token

	if !lr.HasIdentifier() || (!lr.HasEmail() && !lr.HasUsername()) {
		return tokens, ErrInvalidInput
	}
	if !lr.HasPassword() {
		return tokens, ErrInvalidInput
	}

	var user *sqlc.User
	var err error

	switch {
	case lr.HasEmail():
		user, err = a.us.FindOne(ctx, sqlc.CreateUserParams{Email: lr.GetEmail()})
	case lr.HasUsername():
		user, err = a.us.FindOne(ctx, sqlc.CreateUserParams{Username: lr.GetUsername()})
	}

	if err != nil {
		return tokens, fmt.Errorf("database error during user lookup: %w", err)
	}
	if user == nil {
		return tokens, ErrNotFound
	}

	valid, err := a.comparePassword(user.PasswordHash, lr.GetPassword())
	if err != nil {
		return tokens, fmt.Errorf("password comparison failed: %w", err)
	}
	if !valid {
		return tokens, ErrInvalidPassword
	}

	refreshToken, err := generateAndSignToken(
		user.UserID,
		pb.TokenKind_TOKEN_KIND_REFRESH,
		config.GetRefreshTokenDuration(),
		config.GetRefreshTokenSecret(),
		"refresh-token",
	)
	if err != nil {
		return tokens, err
	}
	tokens[0] = refreshToken

	accessToken, err := generateAndSignToken(
		user.UserID,
		pb.TokenKind_TOKEN_KIND_ACCESS,
		config.GetAccessTokenDuration(),
		config.GetAccessTokenSecret(),
		"access-token",
	)
	if err != nil {
		return tokens, err
	}
	tokens[1] = accessToken

	return tokens, nil
}

func (a *userAuth) hashPassword(providedPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(providedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (a *userAuth) comparePassword(hashedPassword, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to compare passwords: %w", err)
	}
	return true, nil
}
