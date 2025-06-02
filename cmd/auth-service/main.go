package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gebhn/auth-service/internal/config"
	"github.com/gebhn/auth-service/internal/db"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	c := db.NewLibsqlConn(config.GetTursoDbUrl(), config.GetTursoDbToken())
	m := db.NewMigrator(c)

	go func() {
		if err := m.Down(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				log.Fatal(err)
			}
		}
		if err := m.Up(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				log.Fatal(err)
			}
		}
	}()

	<-stop
}
