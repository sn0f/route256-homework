package services

import (
	"context"
	"errors"
	"fmt"
	"route256/checkout/internal/domain"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestAddToCart(t *testing.T) {
	type in struct {
		ctx   context.Context
		user  int64
		sku   uint32
		count uint16
	}

	type out struct {
		err error
	}

	var (
		ctx         = context.Background()
		workerCount = 5

		sqlError                = errors.New("sql error")
		serviceUnavailableError = errors.New("service unavailable")
		insufficientStocksError = errors.New("insufficient stocks")
		stocksError             = fmt.Errorf("checking stocks: %v", serviceUnavailableError)
		addToCartRepoError      = fmt.Errorf("adding to cart: %v", sqlError)

		stock = domain.Stock{
			WarehouseID: gofakeit.Int64(),
			Count:       10,
		}
		stocks = []domain.Stock{stock}
	)

	tests := []struct {
		name   string
		in     in
		setup  func(m *mocksService, in *in)
		assert func(require.TestingT, *in, *out)
	}{
		{
			name: "negative case - stocks error",
			in: in{
				ctx:   ctx,
				user:  gofakeit.Int64(),
				sku:   gofakeit.Uint32(),
				count: gofakeit.Uint16(),
			},
			setup: func(m *mocksService, in *in) {
				m.orders.StocksMock.Expect(ctx, in.sku).Return(nil, serviceUnavailableError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.err, stocksError)
			},
		},
		{
			name: "negative case - add to cart repo error",
			in: in{
				ctx:   ctx,
				user:  gofakeit.Int64(),
				sku:   gofakeit.Uint32(),
				count: 5,
			},
			setup: func(m *mocksService, in *in) {
				m.orders.StocksMock.Expect(ctx, in.sku).Return(stocks, nil)
				m.repo.AddToCartMock.Expect(ctx, in.user, in.sku, in.count).Return(sqlError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.err, addToCartRepoError)
			},
		},
		{
			name: "negative case - insufficient stocks error",
			in: in{
				ctx:   ctx,
				user:  gofakeit.Int64(),
				sku:   gofakeit.Uint32(),
				count: 20,
			},
			setup: func(m *mocksService, in *in) {
				m.orders.StocksMock.Expect(ctx, in.sku).Return(stocks, nil)
				m.repo.AddToCartMock.Expect(ctx, in.user, in.sku, in.count).Return(nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.err, insufficientStocksError)
			},
		},
		{
			name: "positive case",
			in: in{
				ctx:   ctx,
				user:  gofakeit.Int64(),
				sku:   gofakeit.Uint32(),
				count: 5,
			},
			setup: func(m *mocksService, in *in) {
				m.orders.StocksMock.Expect(ctx, in.sku).Return(stocks, nil)
				m.repo.AddToCartMock.Expect(ctx, in.user, in.sku, in.count).Return(nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.NoError(t, out.err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ms := newMocksService(t, workerCount)
			tt.setup(ms, &tt.in)
			model := NewModel(ms.orders, ms.products, ms.repo, ms.txManager, ms.workerCount)

			err := model.AddToCart(ctx, tt.in.user, tt.in.sku, tt.in.count)
			tt.assert(t, &tt.in, &out{
				err: err,
			})
		})
	}
}
