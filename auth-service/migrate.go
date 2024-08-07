package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", "postgresql://postgres.lzcfzojsikdlgfbaortw:IPHqieSROTdTNZrv@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?pgbouncer=true")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}

// to create migration file in cli
// migrate create -ext sql -dir db/migrations -seq update_profiles_table
