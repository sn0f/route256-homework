package converters

import (
	"route256/loms/internal/domain"
	api "route256/loms/pkg/loms/v1"
)

func ToStockItemList(items []domain.StockItem) []*api.StockItem {
	res := make([]*api.StockItem, 0, len(items))
	for _, item := range items {
		res = append(res, ToStockItem(item))
	}
	return res
}

func ToStockItem(item domain.StockItem) *api.StockItem {
	return &api.StockItem{
		WarehouseId: item.WarehouseID,
		Count:       uint64(item.Count),
	}
}
