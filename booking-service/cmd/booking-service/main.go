package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/booking-service/db"
	"github.com/reymooy27/arena-backend/booking-service/routes"
)

func main() {

	godotenv.Load(".env")

	db.InitDatabase()
	db.RunMigration()

	router := http.NewServeMux()

	routes.BookingRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("Booking Service is running", "PORT", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}
