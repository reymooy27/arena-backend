package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/reymooy27/arena-backend/auth-service/db"
	"github.com/reymooy27/arena-backend/auth-service/utils"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []User

	query := `SELECT id, username, created_at FROM "user"`
	row, err := db.DB.Query(query)
	if err != nil {
		slog.Error("Query error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	defer row.Close()

	for row.Next() {
		var user User
		if err := row.Scan(&user.Id, &user.Username, &user.CreatedAt); err != nil {
			slog.Error("Query error", err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
			return
		}
		users = append(users, user)
	}

	utils.JSONResponse(w, http.StatusOK, &users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var user User

	query := `SELECT id, username FROM "user" WHERE id = $1 LIMIT 1`

	if err = db.DB.QueryRow(query, parsedId).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Query error", err)
			utils.JSONResponse(w, http.StatusNotFound, "No data with the id")
			return
		}

		slog.Error("Query error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "There is something wrong")
		return
	}

	utils.JSONResponse(w, http.StatusOK, &user)
}
