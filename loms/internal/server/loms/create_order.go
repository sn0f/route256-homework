package loms

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/domain"
	api "route256/loms/pkg/loms/v1"
)

func ValidateCreateOrderRequest(r *api.CreateOrderRequest) error {
	var errs []error

	if r.GetUser() <= 0 {
		errs = append(errs, ErrEmptyUser)
	}
	if len(r.GetItems()) == 0 {
		errs = append(errs, ErrEmptyItems)
	}
	for _, item := range r.Items {
		return ValidateOrderItem(item, errs)
	}

	return errors.Join(errs...)
}

func ValidateOrderItem(r *api.OrderItem, errs []error) error {
	if r.GetSku() == 0 {
		errs = append(errs, ErrEmptySKU)
	}
	if r.GetCount() == 0 {
		errs = append(errs, ErrEmptyCount)
	}

	return errors.Join(errs...)
}

func (s *lomsServer) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	err := ValidateCreateOrderRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	orderItems := ConvertRequestItemsToDomainOrderItems(req.Items)

	orderID, err := s.model.CreateOrder(ctx, req.User, orderItems)
	if err != nil {
		return nil, err
	}

	return &api.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

func ConvertRequestItemsToDomainOrderItems(items []*api.OrderItem) []domain.OrderItem {
	orderItems := make([]domain.OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, domain.OrderItem{
			SKU:   item.Sku,
			Count: uint16(item.Count),
		})
	}
	return orderItems
}
