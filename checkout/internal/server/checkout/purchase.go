package checkout

import (
	"context"
	"fmt"
	api "route256/checkout/pkg/checkout/v1"
)

func ValidatePurchaseRequest(r *api.PurchaseRequest) error {
	if r.GetUser() <= 0 {
		return ErrEmptyUser
	}
	return nil
}

func (s *checkoutServer) Purchase(ctx context.Context, req *api.PurchaseRequest) (*api.PurchaseResponse, error) {
	err := ValidatePurchaseRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	orderID, err := s.model.Purchase(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &api.PurchaseResponse{OrderId: orderID}, nil
}
