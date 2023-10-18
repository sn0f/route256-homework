package loms

import (
	"context"
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/loms/v1"
)

func (c *client) CreateOrder(ctx context.Context, user int64, items []domain.CartItem) (*domain.Order, error) {
	orderItems := make([]*api.OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, &api.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	request := &api.CreateOrderRequest{
		User:  user,
		Items: orderItems,
	}

	response, err := c.lomsClient.CreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	var newOrder = domain.NewOrder(response.OrderId)
	return newOrder, nil
}
