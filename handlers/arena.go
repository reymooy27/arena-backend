package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/reymooy27/arena-backend/db"
)

type BLT struct {
	Id               int8   `json:"id"`
	Nama             string `json:"nama"`
	Pekerjaan        string `json:"pekerjaan"`
	Penghasilan      string `json:"penghasilan"`
	KepemilikanRumah string `json:"kepemilikan_rumah"`
	Aset             string `json:"aset"`
	JumlahTanggungan string `json:"jumlah_tanggungan"`
}

type Data struct {
	Name string
	Age  int
}

type Response struct {
	Message string `json:"message"`
}

func GetArena(w http.ResponseWriter, r *http.Request) {
	var arenas []db.BLT
	result := db.DB.Find(&arenas)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Message: "There was an error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(arenas)
}

func GetArenaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Hello world",
	})
}

func CreateArena(w http.ResponseWriter, r *http.Request) {

	var data Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.Age)

	response := Response{
		Message: "Hello world",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
