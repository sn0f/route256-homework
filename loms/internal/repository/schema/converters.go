package schema

import (
	"route256/loms/internal/domain"
)

func BindSchemaOrderToModelOrder(order Order) domain.Order {
	return domain.Order{
		User:     order.UserID,
		StatusID: domain.OrderStatus(order.StatusID),
	}
}

func BindSchemaStocksToModelStockItems(items []Stock) []domain.StockItem {
	result := make([]domain.StockItem, len(items))
	for i := range items {
		result[i].WarehouseID = items[i].WarehouseID
		result[i].Count = items[i].Count
	}
	return result
}

func BindSchemaReservesToModelReserves(items []Reserve) []domain.Reserve {
	result := make([]domain.Reserve, len(items))
	for i := range items {
		result[i].OrderID = items[i].OrderID
		result[i].WarehouseID = items[i].WarehouseID
		result[i].SKU = uint32(items[i].SKU)
		result[i].Count = items[i].Count
	}
	return result
}

func BindSchemaOrderItemsToModelOrderItems(items []OrderItem) []domain.OrderItem {
	result := make([]domain.OrderItem, len(items))
	for i := range items {
		result[i].SKU = uint32(items[i].SKU)
		result[i].Count = uint16(items[i].Count)
	}
	return result
}

func BindSchemaOrderMessagesToModelOrderMessagess(items []OrderMessage) []domain.OrderMessage {
	result := make([]domain.OrderMessage, len(items))
	for i := range items {
		result[i].ID = items[i].ID
		result[i].OrderID = items[i].OrderID
		result[i].StatusID = domain.OrderStatus(items[i].StatusID)
	}
	return result
}
