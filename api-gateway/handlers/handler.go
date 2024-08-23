package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	pb "github.com/reymooy27/arena-backend/api-gateway/proto/payment"
	"github.com/reymooy27/arena-backend/api-gateway/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Request struct {
	TotalAmount   int    `json:"total_amount"`
	PaymentMethod string `json:"payment_method"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type PaymentHandler struct {
	paymentClient pb.PaymentServiceClient
}

func NewPaymentHandler(apiGateway gateway.APIGateway) *PaymentHandler {
	return &PaymentHandler{paymentClient: apiGateway.PaymentClient}
}

func (s *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var data Request

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Cannot decode json", "err", err)
		utils.JSONResponse(w, 400, fmt.Errorf("Invalid request body"))
		return
	}

	port := os.Getenv("PAYMENT_SERVICE_URL")
	if port == "" {
		slog.Error("PORT not set")
	}

	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Cannot connect to grpc", "err", err)
		utils.JSONResponse(w, 500, fmt.Errorf("Cannot create payment"))
		return
	}

	defer conn.Close()

	req := &pb.PaymentRequest{
		TotalAmount:   int64(data.TotalAmount),
		PaymentMethod: data.PaymentMethod,
	}

	client := pb.NewPaymentServiceClient(conn)

	res, err := client.CreatePayment(context.Background(), req)
	if err != nil {
		slog.Error("Cannot create payment", "err", err)
		utils.JSONResponse(w, 500, fmt.Errorf("Cannot create payment"))
		return
	}

	var response = Response{
		Message: res.Message,
		Success: res.Success,
	}

	utils.JSONResponse(w, 200, response)
}
