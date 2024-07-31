package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/reymooy27/arena-backend/db"
)

type BLT struct {
	Id               int64  `json:"id"`
	Nama             string `json:"nama"`
	Pekerjaan        string `json:"pekerjaan"`
	Penghasilan      string `json:"penghasilan"`
	KepemilikanRumah string `json:"kepemilikan_rumah"`
	Aset             string `json:"aset"`
	JumlahTanggungan string `json:"jumlah_tanggungan"`
}

type Response struct {
	Message string `json:"message"`
}

func GetArena(w http.ResponseWriter, r *http.Request) {
	query := `SELECT * FROM "BLT"`

	row, err := db.DB.Query(query)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("There is something wrong")
	}

	defer row.Close()

	var datas []BLT

	for row.Next() {
		var blt BLT
		if err := row.Scan(&blt.Id, &blt.Aset, &blt.JumlahTanggungan, &blt.KepemilikanRumah, &blt.Nama, &blt.Pekerjaan, &blt.Penghasilan); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("There is something wrong")
		}
		datas = append(datas, blt)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(datas)
}

func GetArenaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("wrong id")
		return
	}

	query := `SELECT * FROM "BLT" WHERE "id" = $1`

	row := db.DB.QueryRow(query, parsedId)

	var blt BLT

	if err := row.Scan(&blt.Id, &blt.Aset, &blt.JumlahTanggungan, &blt.KepemilikanRumah, &blt.Nama, &blt.Pekerjaan, &blt.Penghasilan); err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("No data with the id")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		json.NewEncoder(w).Encode("There is something wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&blt)
}

func CreateArena(w http.ResponseWriter, r *http.Request) {

	var data BLT

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		json.NewEncoder(w).Encode("There is something wrong")
		return
	}

	query := `INSERT INTO "BLT" ("nama","aset","kepemilikan_rumah", "pekerjaan", "penghasilan", "jumlah_tanggungan") VALUES ($1, $2, $3, $4, $5, $6)`
	row, err := db.DB.Exec(query, data.Nama, data.Aset, data.KepemilikanRumah, data.Pekerjaan, data.Penghasilan, data.JumlahTanggungan)

	if err != nil {
		log.Fatal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Cannot insert data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println(row)
	json.NewEncoder(w).Encode("Successfully inserted data")
}

func DeleteArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("wrong id")
		return
	}

	query := `DELETE FROM "BLT" WHERE "id" = $1`
	row, err := db.DB.Exec(query, parsedId)

	if err != nil {
		log.Fatal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Cannot delete data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println(&row)
	json.NewEncoder(w).Encode("Successfully delete data")

}

func UpdateArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("wrong id")
		return
	}

	var data BLT
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		json.NewEncoder(w).Encode("There is something wrong")
		return
	}

	query := `UPDATE "BLT" SET nama = $1,aset = $2, kepemilikan_rumah = $3, pekerjaan = $4, penghasilan = $5, jumlah_tanggungan = $6 WHERE "id" = $7`
	row, err := db.DB.Exec(query, &data.Nama, &data.Aset, &data.KepemilikanRumah, &data.Pekerjaan, &data.Penghasilan, &data.JumlahTanggungan, parsedId)

	if err != nil {
		log.Fatal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Cannot update data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println(&row)
	json.NewEncoder(w).Encode("Successfully update data")

}
