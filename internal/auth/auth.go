package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/config"
)

var (
	ErrInvalidClaims   = errors.New("invalid claims")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNotFound        = errors.New("not found")
	ErrExists          = errors.New("already exists")
)

type Auth interface {
	UserAuth
	TokenAuth
}

type UserAuth interface {
	Register(context.Context, *pb.RegisterRequest) (*pb.User, error)
	Login(context.Context, *pb.LoginRequest) ([2]*pb.Token, error)
}

type TokenAuth interface {
	Verify(context.Context, *pb.Token) (*pb.User, error)
	Refresh(context.Context, *pb.Token) ([2]*pb.Token, error)
	RevokeAccess(context.Context, *pb.Token) (time.Time, error)
	RevokeRefresh(context.Context, *pb.Token) (time.Time, error)
	RevokeAll(context.Context, *pb.User) (time.Time, error)
}

type Claim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var _ UserAuth = (*userAuth)(nil)

func generateAndSignToken(
	userID string,
	kind pb.TokenKind,
	duration time.Duration,
	secret string,
	subject string,
) (*pb.Token, error) {
	expiresAt := time.Now().Add(duration)
	issuedAt := time.Now()
	jti := uuid.New().String()

	claims := &Claim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    config.GetServiceName(),
			Subject:   subject,
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign %s token: %w", subject, err)
	}

	var t pb.Token

	t.SetKind(kind)
	t.SetValue(signedToken)
	t.SetJti(jti)
	t.SetSubject(subject)
	t.SetExpiresAt(timestamppb.New(expiresAt))
	t.SetIssuedAt(timestamppb.New(issuedAt))

	return &t, nil
}
