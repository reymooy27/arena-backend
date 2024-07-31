package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitDatabase() {

	var err error

	// DBConnectionString := "user=postgres.hhqemokocrfgbnaierzf password=Jkhqdlkhqlkwjd235hlqdw9adf! host=aws-0-ap-southeast-1.pooler.supabase.com? port=6543 dbname=postgres"
	DBConnectionString := "postgresql://postgres.hhqemokocrfgbnaierzf:Jkhqdlkhqlkwjd235hlqdw9adf!@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?pgbouncer=true"

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
