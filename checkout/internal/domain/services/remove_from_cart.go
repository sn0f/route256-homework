package services

import (
	"context"
	"fmt"
)

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	return m.transactionManager.RunReadCommitted(ctx, func(ctxTX context.Context) error {
		item, err := m.cartRepository.GetItem(ctx, user, sku)
		if err != nil {
			return fmt.Errorf("getting cart item: %v", err)
		}

		if item.Count > count {
			newCount := item.Count - count
			err = m.cartRepository.UpdateItemCount(ctx, user, sku, newCount)
			if err != nil {
				return fmt.Errorf("updating item count in cart: %v", err)
			}
			return nil
		}

		err = m.cartRepository.DeleteFromCart(ctx, user, sku)
		if err != nil {
			return fmt.Errorf("deleting item from cart: %v", err)
		}
		return nil
	})
}
