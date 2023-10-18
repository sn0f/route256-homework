package domain

import (
	"context"
	"errors"
	"fmt"
)

func (m *Model) CreateOrder(ctx context.Context, user int64, items []OrderItem) (int64, error) {
	orderID, err := m.orderRepository.CreateOrder(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("creating order: %v", err)
	}

	err = m.orderRepository.InsertOrderItems(ctx, orderID, items)
	if err != nil {
		return orderID, fmt.Errorf("inserting order items: %v", err)
	}

	order, err := m.orderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return orderID, fmt.Errorf("getting order: %v", err)
	}

	_, err = m.orderRepository.CreateOrderMessage(ctx, orderID, order.StatusID)
	if err != nil {
		return orderID, fmt.Errorf("inserting order message: %v", err)
	}

	statusID := OrderStatusAwaitingPayment

	var errs []error
	err = m.Reserve(ctx, orderID, items)
	if err != nil {
		errs = append(errs, err)
		statusID = OrderStatusFailed
	}

	err = m.orderRepository.UpdateOrder(ctx, orderID, statusID)
	if err != nil {
		errs = append(errs, fmt.Errorf("updating order status: %v", err))
	}

	m.orderCanceller.AddCancelOrderTask(orderID, m)

	_, err = m.orderRepository.CreateOrderMessage(ctx, orderID, statusID)
	if err != nil {
		errs = append(errs, fmt.Errorf("inserting order message: %v", err))
	}

	return orderID, errors.Join(errs...)
}
