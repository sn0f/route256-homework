package loms

import (
	"context"
	"fmt"
	api "route256/loms/pkg/loms/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func ValidateCancelOrderRequest(r *api.CancelOrderRequest) error {
	if r.GetOrderId() == 0 {
		return ErrEmptyOrder
	}
	return nil
}

func (s *lomsServer) CancelOrder(ctx context.Context, req *api.CancelOrderRequest) (*emptypb.Empty, error) {
	err := ValidateCancelOrderRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	err = s.model.CancelOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
