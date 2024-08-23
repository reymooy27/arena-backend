package routes

import (
	"net/http"

	payment "github.com/reymooy27/arena-backend/api-gateway/handlers"
	arena "github.com/reymooy27/arena-backend/api-gateway/handlers/arena"
	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	"github.com/reymooy27/arena-backend/api-gateway/middleware"
)

func GatewayRoutes(router *http.ServeMux, clients *gateway.APIGateway) {
	router.Handle("POST /payment/create", middleware.AuthMiddleware(http.HandlerFunc(payment.NewPaymentHandler(*clients).CreatePayment)))
	router.Handle("POST /arena/create", middleware.AuthMiddleware(http.HandlerFunc(arena.NewArenaHandler(*clients).CreateArena)))
	router.Handle("GET /arena/{id}", http.HandlerFunc(arena.NewArenaHandler(*clients).GetArenaById))
}
