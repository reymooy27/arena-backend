package main

import (
	// "log"
	// "log/slog"
	// "net/http"
	// "os"

	"context"
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/reymooy27/arena-backend/payment-service/db"
	pb "github.com/reymooy27/arena-backend/payment-service/proto"
	// "github.com/reymooy27/arena-backend/payment-service/routes"
	"google.golang.org/grpc"
)

func main() {

	godotenv.Load(".env")

	db.InitDatabase()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	log.Printf("Server is running at port %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}

}

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}
