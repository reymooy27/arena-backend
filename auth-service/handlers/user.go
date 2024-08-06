package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/reymooy27/arena-backend/auth-service/db"
	"github.com/reymooy27/arena-backend/auth-service/utils"
)

type PublicUser struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot parse query")
		return
	}

	var pageSize int64
	data := r.FormValue("pageSize")
	var err error

	if data != "" {
		pageSize, err = strconv.ParseInt(data, 10, 64)
		if err != nil {
			slog.Error("Parsing query", err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
	} else {
		pageSize = 10
	}

	var users []PublicUser

	query := `SELECT id, username, created_at FROM "user" LIMIT $1`
	row, err := db.DB.Query(query, pageSize)
	if err != nil {
		slog.Error("Query error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	defer row.Close()

	for row.Next() {
		var user PublicUser
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

	var user PublicUser

	query := `SELECT id, username, created_at FROM "user" WHERE id = $1 LIMIT 1`

	if err = db.DB.QueryRow(query, parsedId).Scan(&user.Id, &user.Username, &user.CreatedAt); err != nil {
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

func DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	userData := r.Context().Value("user").(*Claim)

	if int64(userData.Id) != parsedId {
		slog.Error("User not authorized to delete user data")
		utils.JSONResponse(w, http.StatusForbidden, "Not Authorized")
		return
	}

	query := `DELETE FROM "user" WHERE id = $1`
	row, err := db.DB.Exec(query, parsedId)
	if err != nil {
		slog.Error("Cannot delete user", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot delete user")
		return
	}

	rowsAffected, err := row.RowsAffected()

	if err != nil {
		slog.Error("Cannot delete user", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot delete user")
		return
	}

	if rowsAffected == 0 {
		utils.JSONResponse(w, http.StatusNotFound, "No data with the id")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Succesfully delete account")
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	userData := r.Context().Value("user").(*Claim)

	if int64(userData.Id) != parsedId {
		slog.Error("User not authorized to edit user data")
		utils.JSONResponse(w, http.StatusForbidden, "Not Authorized")
		return
	}

	var data User

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := `UPDATE "user" SET username = $1 WHERE id = $2`
	res, err := db.DB.Exec(query, data.Username, parsedId)

	if err != nil {
		slog.Error("Cannot edit user", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot update user")
		return
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		slog.Error("Cannot edit user", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot update user")
		return
	}

	if rowAffected == 0 {
		utils.JSONResponse(w, http.StatusInternalServerError, "No data with id")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Succesfully update account")
}
