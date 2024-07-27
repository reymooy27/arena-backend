package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello sir")
}

func PostData(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  403,
		})
		return
	}

	response := Response{
		Message: "Hello world",
		Status:  200,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
