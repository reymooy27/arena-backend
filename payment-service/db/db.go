package db

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() {

	var err error

	DBConnectionString := os.Getenv("DB_URL")
	if DBConnectionString == "" {
		log.Fatal("DB_URL not set")
	}

	DB, err = sql.Open("postgres", DBConnectionString)

	if err != nil {
		log.Fatal("Error connecting to database: ", "err", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Cannot ping database", err)
	}

	slog.Info("Database connected!")

}
