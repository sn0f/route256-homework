package product

import (
	"context"
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/product/v1"
	"route256/libs/cache"
	"route256/libs/server/rate"
	"time"

	"google.golang.org/grpc"
)

type Client interface {
	GetProduct(ctx context.Context, sku uint32) (*domain.Product, error)
	ListSKUs(ctx context.Context, startAfterSku uint32, count uint32) ([]uint32, error)
}

type GetProductCache interface {
	Get(ctx context.Context, key uint32) (product domain.Product, ok bool)
	Set(ctx context.Context, key uint32, product domain.Product)
}

type client struct {
	productClient   api.ProductServiceClient
	getProductCache GetProductCache
	Token           string
	Limiter         rate.RateLimiter
}

func New(cc *grpc.ClientConn, token string, limiter rate.RateLimiter, ttl time.Duration, size int) *client {
	return &client{
		productClient:   api.NewProductServiceClient(cc),
		getProductCache: cache.New[uint32, domain.Product](size, ttl),
		Token:           token,
		Limiter:         limiter,
	}
}
