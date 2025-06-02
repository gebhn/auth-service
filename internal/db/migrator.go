package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/gebhn/auth-service/internal/config"
)

func NewMigrator(conn *sql.DB) *migrate.Migrate {
	driver, err := sqlite.WithInstance(conn, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(config.GetMigrationDir(), "sqlite", driver)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
