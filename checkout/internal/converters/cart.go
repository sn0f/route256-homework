package converters

import (
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/checkout/v1"
)

func ToListCartResponse(cart domain.Cart) *api.ListCartResponse {
	return &api.ListCartResponse{
		TotalPrice: cart.TotalPrice,
		Items:      ToCartItemList(cart.Items),
	}
}

func ToCartItemList(items []domain.CartItem) []*api.CartItem {
	res := make([]*api.CartItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToCartItem(item))
	}
	return res
}

func ToCartItem(item domain.CartItem) *api.CartItem {
	return &api.CartItem{
		Sku:   item.SKU,
		Count: uint32(item.Count),
		Name:  item.Name,
		Price: item.Price,
	}
}
