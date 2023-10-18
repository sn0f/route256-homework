package domain

import (
	"context"
	"errors"
	"fmt"
)

// Резервируем товары, но не списываем стоки до оплаты.
// Поэтому из стоков вычитаем резервы по неоплаченным заявкам.
func (m *Model) Reserve(ctx context.Context, orderID int64, items []OrderItem) error {
	return m.transactionManager.RunSerializable(ctx, func(ctx context.Context) error {
		for _, item := range items {
			stocks, err := m.orderRepository.GetStocks(ctx, item.SKU)
			if err != nil {
				return fmt.Errorf("getting stocks: %v", err)
			}

			remainingCount := uint64(item.Count)

			for _, stock := range stocks {
				reservedCount, err := m.orderRepository.GetReservedCount(ctx, stock.WarehouseID, item.SKU)
				if err != nil {
					return fmt.Errorf("getting reserved count: %v", err)
				}

				reserve := Reserve{OrderID: orderID, WarehouseID: stock.WarehouseID, SKU: item.SKU}
				stockCount := stock.Count - reservedCount

				if stockCount <= 0 {
					continue
				}

				if remainingCount > stockCount {
					remainingCount -= stockCount
					reserve.Count = stockCount
				} else {
					reserve.Count = remainingCount
					remainingCount = 0
				}

				err = m.orderRepository.CreateReserve(ctx, reserve)
				if err != nil {
					return fmt.Errorf("creating reserve: %v", err)
				}

				if remainingCount == 0 {
					break
				}
			}

			if remainingCount > 0 {
				return errors.New("insufficient stocks")
			}
		}
		return nil
	})
}
