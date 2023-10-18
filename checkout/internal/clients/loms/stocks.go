package loms

import (
	"context"
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/loms/v1"
)

func (c *client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	request := &api.StocksRequest{Sku: sku}

	response, err := c.lomsClient.Stocks(ctx, request)
	if err != nil {
		return nil, err
	}

	stocks := make([]domain.Stock, 0, len(response.Stocks))
	for _, stock := range response.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.WarehouseId,
			Count:       stock.Count,
		})
	}

	return stocks, nil
}
