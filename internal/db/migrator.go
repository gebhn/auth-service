package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/gebhn/auth-service/build/package/auth-service/migrations"
)

func NewMigrator(conn *sql.DB) *migrate.Migrate {
	driver, err := sqlite.WithInstance(conn, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite", driver)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
