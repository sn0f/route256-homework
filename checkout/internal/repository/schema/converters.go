package schema

import (
	"route256/checkout/internal/domain"
)

func BindSchemaCartItemsToModelCartItems(items []CartItem) []domain.CartItem {
	result := make([]domain.CartItem, len(items))
	for i := range items {
		result[i].SKU = items[i].SKU
		result[i].Count = items[i].Count
	}
	return result
}
