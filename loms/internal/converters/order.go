package converters

import (
	"route256/loms/internal/domain"
	api "route256/loms/pkg/loms/v1"
)

func ToListOrderResponse(order domain.Order) *api.ListOrderResponse {
	return &api.ListOrderResponse{
		User:   order.User,
		Status: ToResponseOrderStatus(order.StatusID),
		Items:  ToOrderItemList(order.Items),
	}
}

func ToOrderItemList(items []domain.OrderItem) []*api.OrderItem {
	res := make([]*api.OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToOrderItem(item))
	}
	return res
}

func ToOrderItem(item domain.OrderItem) *api.OrderItem {
	return &api.OrderItem{
		Sku:   item.SKU,
		Count: uint32(item.Count),
	}
}

func ToResponseOrderStatus(s domain.OrderStatus) api.OrderStatus {
	switch s {
	case domain.OrderStatusNew:
		return 1
	case domain.OrderStatusAwaitingPayment:
		return 2
	case domain.OrderStatusFailed:
		return 3
	case domain.OrderStatusPayed:
		return 4
	case domain.OrderStatusCancelled:
		return 5
	}
	return 0
}
