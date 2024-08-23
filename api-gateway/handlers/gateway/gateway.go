package gateway

import (
	pbArena "github.com/reymooy27/arena-backend/api-gateway/proto/arena"
	pbPayment "github.com/reymooy27/arena-backend/api-gateway/proto/payment"
	"google.golang.org/grpc"
)

type APIGateway struct {
	ArenaClient   pbArena.ArenaServiceClient
	PaymentClient pbPayment.PaymentServiceClient
}

func NewAPIGateway(arenaConn, paymentConn *grpc.ClientConn) *APIGateway {
	return &APIGateway{
		ArenaClient:   pbArena.NewArenaServiceClient(arenaConn),
		PaymentClient: pbPayment.NewPaymentServiceClient(paymentConn),
	}
}
