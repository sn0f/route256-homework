package domain

import (
	"context"
	"fmt"
)

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]StockItem, error) {
	items, err := m.orderRepository.GetStocks(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("getting stocks: %v", err)
	}

	return items, nil
}
