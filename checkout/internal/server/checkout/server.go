package checkout

import (
	"route256/checkout/internal/domain/services"
	api "route256/checkout/pkg/checkout/v1"
)

type checkoutServer struct {
	api.UnimplementedCheckoutServiceServer
	model *services.Model
}

func NewCheckoutServer(model *services.Model) *checkoutServer {
	return &checkoutServer{
		model: model,
	}
}
