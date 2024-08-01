package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/db"
	"github.com/reymooy27/arena-backend/routes"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db.InitDatabase()

	router := http.NewServeMux()

	routes.ArenaRoutes(router)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Println("Server running on port 8000")
	server.ListenAndServe()
}
