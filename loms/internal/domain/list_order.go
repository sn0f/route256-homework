package domain

import (
	"context"
	"fmt"
)

func (m *Model) ListOrder(ctx context.Context, orderID int64) (Order, error) {
	order, err := m.orderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("getting order: %v", err)
	}

	order.Items, err = m.orderRepository.GetOrderItems(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("getting order items: %v", err)
	}

	return order, nil
}
