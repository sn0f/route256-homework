package domain

import (
	"context"
	"errors"
	"fmt"
)

func (m *Model) OrderPayed(ctx context.Context, orderID int64) error {
	txErr := m.transactionManager.RunRepeatableRead(ctx, func(ctx context.Context) error {
		order, err := m.orderRepository.GetOrder(ctx, orderID)
		if err != nil {
			return fmt.Errorf("getting order: %v", err)
		}

		if order.StatusID != OrderStatusAwaitingPayment {
			return errors.New("wrong order status")
		}

		status := OrderStatusPayed
		err = m.orderRepository.UpdateOrder(ctx, orderID, OrderStatusPayed)
		if err != nil {
			return fmt.Errorf("paying order: %v", err)
		}

		_, err = m.orderRepository.CreateOrderMessage(ctx, orderID, status)
		if err != nil {
			return fmt.Errorf("inserting order message: %v", err)
		}

		reserves, err := m.orderRepository.GetReserves(ctx, orderID)
		if err != nil {
			return fmt.Errorf("getting reserves by order: %v", err)
		}

		// Списываем стоки после оплаты
		for _, reserve := range reserves {
			err = m.orderRepository.SubtractStocks(ctx, reserve.WarehouseID, reserve.SKU, reserve.Count)
			if err != nil {
				return fmt.Errorf("subtracking stocks: %v", err)
			}
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	m.orderCanceller.RemoveCancelOrderTask(orderID)
	return nil
}
