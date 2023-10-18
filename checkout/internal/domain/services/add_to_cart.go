package services

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := m.orderManager.Stocks(ctx, sku)
	if err != nil {
		return fmt.Errorf("checking stocks: %v", err)
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			err = m.cartRepository.AddToCart(ctx, user, sku, count)
			if err != nil {
				return fmt.Errorf("adding to cart: %v", err)
			}
			return nil
		}
	}

	return ErrInsufficientStocks
}
