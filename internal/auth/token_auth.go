package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/config"
	"github.com/gebhn/auth-service/internal/db/sqlc"
	"github.com/gebhn/auth-service/internal/store"
)

type tokenAuth struct {
	us *store.UserStore
	ts *store.TokenStore
}

func NewTokenAuth(us *store.UserStore, ts *store.TokenStore) *tokenAuth {
	return &tokenAuth{
		us: us,
		ts: ts,
	}
}

func (t *tokenAuth) Verify(ctx context.Context, tkn *pb.Token) (*pb.User, error) {
	if !tkn.HasValue() {
		return nil, ErrInvalidInput
	}
	claims, err := t.verifyToken(tkn.GetValue(), config.GetAccessTokenSecret())
	if err != nil {
		return nil, err
	}
	user, err := t.us.FindOne(ctx, sqlc.CreateUserParams{UserID: claims.UserID})
	if err != nil {
		return nil, err
	}

	var u pb.User
	u.SetUserId(user.UserID)
	u.SetEmail(user.Email)
	u.SetUsername(user.Username)

	return &u, nil
}

func (t *tokenAuth) Refresh(ctx context.Context, tkn *pb.Token) ([2]*pb.Token, error) {
	var tokens [2]*pb.Token

	if !tkn.HasValue() {
		return tokens, ErrInvalidInput
	}
	claims, err := t.verifyToken(tkn.GetValue(), config.GetRefreshTokenSecret())
	if err != nil {
		return tokens, err
	}

	// TODO @gebhn: Invalidate outstanding refresh token

	refreshToken, err := generateAndSignToken(
		claims.UserID,
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
		claims.UserID,
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

// func (t *tokenAuth) RevokeAccess(context.Context, *pb.Token) (time.Time, error) {
// }
//
// func (t *tokenAuth) RevokeRefresh(context.Context, *pb.Token) (time.Time, error) {
// }
//
// func (t *tokenAuth) RevokeAll(context.Context, *pb.User) (time.Time, error) {
// }

func (t *tokenAuth) verifyToken(token string, secret string) (*Claim, error) {
	tkn, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tkn.Claims.(*Claim)
	if !ok || !tkn.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}
