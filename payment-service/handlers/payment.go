package handlers

import (
	// "encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
	"google.golang.org/protobuf/proto"

	// "github.com/reymooy27/arena-backend/payment-service/db"
	pb "github.com/reymooy27/arena-backend/payment-service/proto"
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

	// user := r.Context().Value("user").(*Claim)

	req := &pb.HelloRequest{}

	br, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Read body error", "message", err)
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := proto.Unmarshal(br, req); err != nil {
		slog.Error("Unmarshal body error", "message", err)
		utils.JSONResponse(w, 500, "Invalid request body")
		return
	}

	name := req.GetName()
	log.Printf(req.GetName())

	res := &pb.HelloResponse{Message: "Hello " + name}

	response, err := proto.Marshal(res)
	if err != nil {
		slog.Error("Marshal response error", "message", err)
		utils.JSONResponse(w, 500, "Cannot marshal response")
		return
	}

	w.Write(response)

	// query := `INSERT INTO payments (user_id, amount) VALUES ($1, $2)`
	// result, err := db.DB.Exec(query, user.Id, body.Amount)
	// log.Println(result)
	// if err != nil {
	// 	slog.Error("Query error", "message", err)
	// 	utils.JSONResponse(w, http.StatusBadRequest, "Cannot create payment")
	// 	return
	// }

	// utils.JSONResponse(w, http.StatusOK, "Payment succesfull")
}
