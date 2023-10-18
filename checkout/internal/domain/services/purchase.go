package services

import (
	"context"
	"fmt"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {
	cart, err := m.ListCart(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("list cart: %v", err)
	}

	order, err := m.orderManager.CreateOrder(ctx, user, cart.Items)
	if err != nil {
		return 0, fmt.Errorf("creating order: %v", err)
	}

	logger.Info("order created", zap.String("order", fmt.Sprintf("%+v", *order)))

	err = m.cartRepository.ClearCart(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("clearing cart: %v", err)
	}

	return order.OrderID, nil
}
