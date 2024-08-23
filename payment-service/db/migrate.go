package db

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigration() {

	connString := os.Getenv("DB_URL")
	if connString == "" {
		slog.Error("DB_URL not set")
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Failed to create migration driver", "err", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		slog.Error("Failed to create migrate instance", "err", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to run migrations", "err", err)
	}

	v, d, err := m.Version()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	slog.Info("Migration", "Version", v, "Dirty", d)

	slog.Info("Migrations applied successfully!")
}

// to create migration file in cli
// migrate create -ext sql -dir db/migrations -seq update_profiles_table
