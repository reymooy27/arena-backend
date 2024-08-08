package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/booking-service/handlers"
	"github.com/reymooy27/arena-backend/booking-service/middleware"
)

func BookingRoutes(router *http.ServeMux) {
	router.Handle("GET /bookings", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserBookings)))
	router.Handle("POST /booking", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateBooking)))
}
