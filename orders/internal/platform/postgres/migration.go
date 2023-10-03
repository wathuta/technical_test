package database

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/wathuta/technical_test/orders/internal/config"
	"golang.org/x/exp/slog"
)

func RunMigrations() error {
	connUri, _ := config.ConnectionURLBuilder("postgres-migrate")
	m, err := migrate.New(
		"file://internal/platform/migrations",
		*connUri,
	)

	if err != nil {
		slog.Error("migration connection failed", "error", err)
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migration up failed", "error", err, m.Log.Verbose())
		return err
	}

	log.Println("migrations finished successfully")
	return nil
}
