package main

import (
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

	router := http.NewServeMux()

	routes.BookingRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info(`Server running on`, "PORT", port)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error(`Server error`, "message", err)
	}
}
