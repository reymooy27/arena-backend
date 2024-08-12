package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/payment-service/handlers"
	// "github.com/reymooy27/arena-backend/payment-service/middleware"
)

func PaymentRoutes(router *http.ServeMux) {
	// router.Handle("POST /payment", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePayment)))
	router.HandleFunc("POST /payment", handlers.CreatePayment)
}
