package routes

import (
	"net/http"

	arena "github.com/reymooy27/arena-backend/api-gateway/handlers/arena"
	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	payment "github.com/reymooy27/arena-backend/api-gateway/handlers/payment"
	"github.com/reymooy27/arena-backend/api-gateway/middleware"
)

func GatewayRoutes(router *http.ServeMux, clients *gateway.APIGateway) {
	//payment
	router.Handle("POST /payment/create", middleware.AuthMiddleware(http.HandlerFunc(payment.NewPaymentHandler(*clients).CreatePayment)))

	//arena
	router.Handle("POST /arena/create", middleware.AuthMiddleware(http.HandlerFunc(arena.NewArenaHandler(*clients).CreateArena)))
	router.Handle("GET /arena/{id}", http.HandlerFunc(arena.NewArenaHandler(*clients).GetArenaById))
	router.Handle("DELETE /arena/{id}", http.HandlerFunc(arena.NewArenaHandler(*clients).DeleteArena))
	router.Handle("PUT /arena/{id}", http.HandlerFunc(arena.NewArenaHandler(*clients).UpdateArena))
	router.Handle("GET /arenas", http.HandlerFunc(arena.NewArenaHandler(*clients).GetArenas))
}
