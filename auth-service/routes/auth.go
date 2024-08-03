package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/auth-service/handlers"
	"github.com/reymooy27/arena-backend/auth-service/middleware"
)

func AuthRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", handlers.Login)
	router.HandleFunc("POST /signup", handlers.Signup)
	router.HandleFunc("GET /logout", handlers.Logout)
	router.Handle("GET /verify", middleware.AuthMiddleware(http.HandlerFunc(handlers.Verify)))
}
