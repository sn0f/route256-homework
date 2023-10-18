package loms

import (
	"context"
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/loms/v1"

	"google.golang.org/grpc"
)

type Client interface {
	CreateOrder(ctx context.Context, user int64, items []domain.CartItem) (*domain.Order, error)
	Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error)
}

type client struct {
	lomsClient api.LomsServiceClient
}

func New(cc *grpc.ClientConn) *client {
	return &client{
		lomsClient: api.NewLomsServiceClient(cc),
	}
}
