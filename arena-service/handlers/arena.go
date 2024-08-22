package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/reymooy27/arena-backend/arena-service/db"
	"github.com/reymooy27/arena-backend/arena-service/utils"
)

type Arena struct {
	Id          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func GetArena(w http.ResponseWriter, r *http.Request) {
	query := `SELECT * FROM "arenas"`

	row, err := db.DB.Query(query)

	if err != nil {
		slog.Error("Query arena", "err", err)
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid query")
	}

	defer row.Close()

	var datas []Arena

	for row.Next() {
		var blt Arena
		if err := row.Scan(&blt.Id, &blt.CreatedAt, &blt.Name, &blt.Description); err != nil {
			slog.Error("Scan arena", "err", err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
			return
		}
		datas = append(datas, blt)
	}

	utils.JSONResponse(w, http.StatusOK, datas)
}

func GetArenaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {

		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	query := `SELECT * FROM "arenas" WHERE "id" = $1`

	row := db.DB.QueryRow(query, parsedId)

	var blt Arena

	if err := row.Scan(&blt.Id, &blt.CreatedAt, &blt.Name, &blt.Description); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Scan arena", "err", err)
			utils.JSONResponse(w, http.StatusNotFound, "No data with the id")
			return
		}

		utils.JSONResponse(w, http.StatusInternalServerError, "There is something wrong")
		return
	}

	utils.JSONResponse(w, http.StatusOK, &blt)
}

func CreateArena(w http.ResponseWriter, r *http.Request) {

	var data Arena

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := `INSERT INTO "arenas" ("name","description") VALUES ($1, $2)`
	_, err = db.DB.Exec(query, data.Name, data.Description)

	if err != nil {
		slog.Error("Query insert arena", "err", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot insert data")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Successfully inserted data")
}

func DeleteArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	query := `DELETE FROM "arenas" WHERE "id" = $1`
	row, err := db.DB.Exec(query, parsedId)

	rowsAffected, err := row.RowsAffected()

	if err != nil {
		slog.Error("Cannot delete arena", "err", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Error checking rows affected")
		return
	}

	if rowsAffected == 0 {
		utils.JSONResponse(w, http.StatusNotFound, "No data with the id")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Successfully delete data")

}

func UpdateArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var data Arena
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Invalid request body")
		return
	}

	query := `UPDATE "arenas" SET name = $1, description = $2 WHERE "id" = $3`
	result, err := db.DB.Exec(query, &data.Name, &data.Description, parsedId)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		slog.Error("Cannot update arena", "err", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Error checking rows affected")
		return
	}

	if rowsAffected == 0 {
		utils.JSONResponse(w, http.StatusNotFound, "No data with the id")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Successfully update data")

}
