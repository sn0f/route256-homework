package checkout

import (
	"context"
	"fmt"
	"route256/checkout/internal/converters"
	api "route256/checkout/pkg/checkout/v1"
)

func ValidateListCartRequest(r *api.ListCartRequest) error {
	if r.GetUser() <= 0 {
		return ErrEmptyUser
	}
	return nil
}

func (s *checkoutServer) ListCart(ctx context.Context, req *api.ListCartRequest) (*api.ListCartResponse, error) {
	err := ValidateListCartRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	cart, err := s.model.ListCart(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return converters.ToListCartResponse(cart), nil
}
