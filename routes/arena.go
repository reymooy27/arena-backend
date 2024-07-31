package routes

import (
	"github.com/reymooy27/arena-backend/handlers"
	"net/http"
)

func ArenaRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /arena", handlers.GetArena)
	router.HandleFunc("GET /arena/{id}", handlers.GetArenaById)
	router.HandleFunc("POST /arena", handlers.CreateArena)
	router.HandleFunc("DELETE /arena/{id}", handlers.DeleteArena)
	router.HandleFunc("PUT /arena/{id}", handlers.UpdateArena)
}
