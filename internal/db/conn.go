package db

import (
	"database/sql"
	"log"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewLibsqlConn(connectionString, authToken string) *sql.DB {
	conn, err := sql.Open("libsql", connectionString+"?authToken="+authToken)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
