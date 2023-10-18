package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/converters"
	api "route256/loms/pkg/loms/v1"
)

func ValidateListOrderRequest(r *api.ListOrderRequest) error {
	if r.GetOrderId() == 0 {
		return ErrEmptyOrder
	}
	return nil
}

func (s *lomsServer) ListOrder(ctx context.Context, req *api.ListOrderRequest) (*api.ListOrderResponse, error) {
	err := ValidateListOrderRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	order, err := s.model.ListOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return converters.ToListOrderResponse(order), nil
}
