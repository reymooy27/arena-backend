package main

import (
	"log/slog"
	"net/http"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("Server running on port 8001")
	server.ListenAndServe()
}
