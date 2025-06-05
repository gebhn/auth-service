package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gebhn/auth-service/api/pb"
)

var kinds = [5]time.Duration{
	pb.TokenKind_TOKEN_KIND_UNKNOWN:            0,
	pb.TokenKind_TOKEN_KIND_REFRESH:            getRefreshTokenDuration(),
	pb.TokenKind_TOKEN_KIND_ACCESS:             getAccessTokenDuration(),
	pb.TokenKind_TOKEN_KIND_PASSWORD_RESET:     0,
	pb.TokenKind_TOKEN_KIND_EMAIL_VERIFICATION: 0,
}

func GetTursoDbUrl() string {
	return readEnvVar("TURSO_DB_URL", "libsql://my-super-db.turso.io")
}

func GetTursoDbToken() string {
	return readEnvVar("TURSO_DB_TOKEN", "super.secret_token")
}

func GetRedisAddress() string {
	return readEnvVar("REDIS_ADDRESS", "localhost:6379")
}

func GetRedisPassword() string {
	return readEnvVar("REDIS_PASSWORD", "password")
}

func GetGrpcServerPort() string {
	return readEnvVar("GRPC_SERVER_PORT", "50051")
}

func GetServiceName() string {
	return readEnvVar("SERVICE_NAME", "auth-service-1")
}

func GetRefreshTokenSecret() string {
	return readEnvVar("REFRESH_TOKEN_SECRET", "keep-it-secret-keep-it-safe")
}

func GetAccessTokenSecret() string {
	return readEnvVar("ACCESS_TOKEN_SECRET", "is-it-secret-?-is-it-safe-?")
}

func GetTokenDuration(kind pb.TokenKind) time.Duration {
	return kinds[kind]
}

func getRefreshTokenDuration() time.Duration {
	return time.Hour * 24 * 7
}

func getAccessTokenDuration() time.Duration {
	return time.Minute * 5
}

func readEnvVar(envVar, suggestion string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	panic(fmt.Sprintf("env var %s is not set, suggested value: %s", envVar, suggestion))
}
