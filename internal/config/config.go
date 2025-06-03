package config

import (
	"fmt"
	"os"
	"time"
)

func GetTursoDbUrl() string {
	return readEnvVar("TURSO_DB_URL", "libsql://my-super-db.turso.io")
}

func GetTursoDbToken() string {
	return readEnvVar("TURSO_DB_TOKEN", "super.secret_token")
}

func GetGrpcServerPort() string {
	return readEnvVar("GRPC_SERVER_PORT", "50051")
}

func GetMigrationDir() string {
	return readEnvVar("MIGRATION_DIR", "file:///build/package/service/migrations/")
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

func GetRefreshTokenDuration() time.Duration {
	return time.Hour * 24 * 7
}

func GetAccessTokenDuration() time.Duration {
	return time.Minute * 5
}

func readEnvVar(envVar, suggestion string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	panic(fmt.Sprintf("env var %s is not set, suggested value: %s", envVar, suggestion))
}
