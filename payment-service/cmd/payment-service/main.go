package main

import (
	"log/slog"
	"net"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/payment-service/db"
	"github.com/reymooy27/arena-backend/payment-service/internal/service/payment"
	pb "github.com/reymooy27/arena-backend/payment-service/proto"
	"google.golang.org/grpc"
)

func main() {

	godotenv.Load(".env")

	db.InitDatabase()
	db.RunMigration()

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		slog.Error("Could not start grpc server", "err", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterPaymentServiceServer(s, &payment.Server{})

	slog.Info("Payment GRPC is running at", "PORT", listener.Addr())

	if err := s.Serve(listener); err != nil {
		slog.Error("Could not start grpc server", "err", err)
	}

}
