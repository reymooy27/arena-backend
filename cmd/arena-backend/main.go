package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/db"
	"github.com/reymooy27/arena-backend/routes"
)

func main() {

	godotenv.Load(".env")

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
