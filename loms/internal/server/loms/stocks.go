package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/converters"
	api "route256/loms/pkg/loms/v1"
)

func ValidateStocksRequest(r *api.StocksRequest) error {
	if r.GetSku() == 0 {
		return ErrEmptySKU
	}
	return nil
}

func (s *lomsServer) Stocks(ctx context.Context, req *api.StocksRequest) (*api.StocksResponse, error) {
	err := ValidateStocksRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	items, err := s.model.Stocks(ctx, req.Sku)
	if err != nil {
		return nil, err
	}

	return &api.StocksResponse{Stocks: converters.ToStockItemList(items)}, nil
}
