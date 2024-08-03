package routes

import (
	"net/http"

	"github.com/reymooy27/arena-backend/arena-service/handlers"
	"github.com/reymooy27/arena-backend/arena-service/middleware"
)

func ArenaRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /arena", handlers.GetArena)
	router.HandleFunc("GET /arena/{id}", handlers.GetArenaById)
	router.HandleFunc("POST /arena", handlers.CreateArena)
	router.Handle("DELETE /arena/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteArena)))
	router.HandleFunc("PUT /arena/{id}", handlers.UpdateArena)
}
