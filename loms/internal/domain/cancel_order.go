package domain

import (
	"context"
	"errors"
	"fmt"
)

func (m *Model) CancelOrder(ctx context.Context, orderID int64) error {
	txErr := m.transactionManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		order, err := m.orderRepository.GetOrder(ctx, orderID)
		if err != nil {
			return fmt.Errorf("getting order: %v", err)
		}

		if order.StatusID == OrderStatusPayed || order.StatusID == OrderStatusCancelled {
			return errors.New("wrong order status")
		}

		err = m.orderRepository.DeleteReserves(ctx, orderID)
		if err != nil {
			return fmt.Errorf("deleting reserves: %v", err)
		}

		status := OrderStatusCancelled
		err = m.orderRepository.UpdateOrder(ctx, orderID, status)
		if err != nil {
			return fmt.Errorf("updating order status: %v", err)
		}

		_, err = m.orderRepository.CreateOrderMessage(ctx, orderID, status)
		if err != nil {
			return fmt.Errorf("inserting order message: %v", err)
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	m.orderCanceller.RemoveCancelOrderTask(orderID)
	return nil
}
