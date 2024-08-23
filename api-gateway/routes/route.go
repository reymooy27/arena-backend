package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/api-gateway/handlers"
	"github.com/reymooy27/arena-backend/api-gateway/middleware"
)

func GatewayRoutes(router *http.ServeMux) {
	router.Handle("POST /payment/create", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePayment)))
}
