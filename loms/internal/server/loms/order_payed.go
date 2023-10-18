package loms

import (
	"context"
	"fmt"

	api "route256/loms/pkg/loms/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func ValidateOrderPayed(req *api.OrderPayedRequest) error {
	if req.GetOrderId() == 0 {
		return ErrEmptyOrder
	}
	return nil
}

func (s *lomsServer) OrderPayed(ctx context.Context, req *api.OrderPayedRequest) (*emptypb.Empty, error) {
	err := ValidateOrderPayed(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	err = s.model.OrderPayed(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
