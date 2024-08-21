package main

import (
	"log"
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
	db.RunMigration()

	router := http.NewServeMux()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("Auth Service is running", "PORT", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}
