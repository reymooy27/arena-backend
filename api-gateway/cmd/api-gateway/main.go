package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	"github.com/reymooy27/arena-backend/api-gateway/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	godotenv.Load(".env")

	router := http.NewServeMux()

	ARENA_SERVICE_URL := os.Getenv("ARENA_SERVICE_URL")
	PAYMENT_SERVICE_URL := os.Getenv("PAYMENT_SERVICE_URL")

	arenaConn, err := grpc.NewClient(ARENA_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Cannot connect to arena grpc", "err", err)
		return
	}

	defer arenaConn.Close()

	paymentConn, err := grpc.NewClient(PAYMENT_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Cannot connect to payment grpc", "err", err)
		return
	}

	defer paymentConn.Close()

	apiGateway := gateway.NewAPIGateway(arenaConn, paymentConn)

	routes.GatewayRoutes(router, apiGateway)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("API Gateway is running", "PORT", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot start API Gateway: %s", err)
	}
}
