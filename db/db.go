package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() {

	var err error

	DBConnectionString := os.Getenv("DB_URL")
	DB, err = sql.Open("postgres", DBConnectionString)

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Cannot ping database")
	}

	fmt.Println("Database connected!")

}
