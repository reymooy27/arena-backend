package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/reymooy27/arena-backend/payment-service/db"
	"github.com/reymooy27/arena-backend/payment-service/utils"
)

type Body struct {
	Amount int `json:"amount"`
}

type Claim struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// INFO: not finished,
// INFO: still testing
func CreatePayment(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*Claim)

	var body Body

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := `INSERT INTO payments (user_id, amount) VALUES ($1, $2)`
	result, err := db.DB.Exec(query, user.Id, body.Amount)
	log.Println(result)
	if err != nil {
		slog.Error("Query error", "message", err)
		utils.JSONResponse(w, http.StatusBadRequest, "Cannot create payment")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Payment succesfull")
}
