package payment

import (
	"context"
	"log/slog"

	"github.com/reymooy27/arena-backend/payment-service/db"
	pb "github.com/reymooy27/arena-backend/payment-service/proto"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *Server) CreatePayment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	user_id := 1
	booking_id := 1

	var last_inserted_id int64

	query := `INSERT INTO payments (user_id, total_amount, payment_method, booking_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.DB.QueryRow(query, user_id, in.TotalAmount, in.PaymentMethod, booking_id).Scan(&last_inserted_id)

	if err != nil {
		slog.Error("Could not insert data", "err", err)
		return &pb.PaymentResponse{Message: "Payment failed", Success: false}, nil
	}

	return &pb.PaymentResponse{Message: "Payment success", Success: true, PaymentId: last_inserted_id}, nil
}

func (s *Server) CancelPayment(ctx context.Context, in *pb.CancelPaymentRequest) (*pb.CancelPaymentResponse, error) {
	query := `DELETE FROM payments WHERE id = $1`

	res, err := db.DB.ExecContext(ctx, query, in.PaymentId)

	if err != nil {
		slog.Error("Could not insert data", "err", err)
		return &pb.CancelPaymentResponse{Message: "Cancel payment failed", Success: false}, nil
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		slog.Error("Payment not found")
		return &pb.CancelPaymentResponse{Message: "Payment data not found", Success: false}, nil
	}

	return &pb.CancelPaymentResponse{Message: "Cancel payment success", Success: true}, nil
}
