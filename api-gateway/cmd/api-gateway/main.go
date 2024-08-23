package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/api-gateway/routes"
)

func main() {

	godotenv.Load(".env")

	router := http.NewServeMux()

	routes.GatewayRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("API Gateway is running", "PORT", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot start API Gateway: %s", err)
	}
}
