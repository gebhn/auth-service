package db

import (
	"database/sql"
	"log"

	"github.com/tursodatabase/libsql-client-go/libsql"
)

func NewLibsqlConn(connectionString, authToken string) *sql.DB {
	opts := []libsql.Option{}
	if authToken != "" {
		opts = append(opts, libsql.WithAuthToken(authToken))
	}

	driver, err := libsql.NewConnector(connectionString, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return sql.OpenDB(driver)
}
