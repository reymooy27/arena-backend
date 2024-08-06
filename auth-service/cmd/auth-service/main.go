package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/auth-service/db"
	"github.com/reymooy27/arena-backend/auth-service/routes"
)

func main() {

	godotenv.Load(".env")

	db.InitDatabase()

	router := http.NewServeMux()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	server := http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	log.Println("Server running on port 8001")
	server.ListenAndServe()
}
