package auth

import (
	"context"
	"time"

	"github.com/gebhn/auth-service/api/pb"
)

type Auth interface {
	UserAuth
	TokenAuth
}

type UserAuth interface {
	Register(context.Context) (*pb.User, error)
	Login(context.Context) ([2]*pb.Token, error)
}

type TokenAuth interface {
	Verify(context.Context) (*pb.User, error)
	Refresh(context.Context, *pb.Token) ([2]*pb.Token, error)
	RevokeAccess(context.Context, *pb.Token) (time.Time, error)
	RevokeRefresh(context.Context, *pb.Token) (time.Time, error)
	RevokeAll(context.Context, *pb.User) (time.Time, error)
}
