package gateway

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	pb "github.com/reymooy27/arena-backend/payment-service/proto"
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

func APIGateway() {
	router := http.NewServeMux()
	router.HandleFunc("POST /payment/create", createPayment)
	router.HandleFunc("POST /paymet/cancel", cancelPayment)

	port := "5000"

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("Starting server", "port", port)

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Cannot start server", "err", err)
		return
	}
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	var data Request

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Cannot decode json", "err", err)
		return
	}

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Cannot connect to grpc", "err", err)
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
		return
	}

	var response = Response{
		Message: res.Message,
		Success: res.Success,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&response)
}

func cancelPayment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
