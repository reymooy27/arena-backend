package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/auth-service/handlers"
	"github.com/reymooy27/arena-backend/auth-service/middleware"
)

func UserRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /user/{id}", handlers.GetUserByID)
	router.HandleFunc("GET /users", handlers.GetUsers)
	router.Handle("DELETE /user/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteUserAccount)))
	router.Handle("PUT /user/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.EditUser)))
}
