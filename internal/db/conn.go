package db

import (
	"database/sql"
	"log"

	"github.com/tursodatabase/libsql-client-go/libsql"
)

func NewLibsqlConn(connectionString, authToken string) *sql.DB {
	driver, err := libsql.NewConnector(connectionString, libsql.WithAuthToken(authToken))
	if err != nil {
		log.Fatal(err)
	}

	return sql.OpenDB(driver)
}
