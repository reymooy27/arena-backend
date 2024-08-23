package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/arena-service/db"
	"github.com/reymooy27/arena-backend/arena-service/internal/service/arena"
	pb "github.com/reymooy27/arena-backend/arena-service/proto"
	"google.golang.org/grpc"
)

func main() {

	godotenv.Load(".env")

	db.InitDatabase()
	db.RunMigration()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		slog.Error("Could not start grpc server", "err", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterArenaServiceServer(s, &arena.Server{})

	slog.Info("Arena Service is running at", "PORT", listener.Addr())

	if err := s.Serve(listener); err != nil {
		slog.Error("Could not start grpc server", "err", err)
	}
}
