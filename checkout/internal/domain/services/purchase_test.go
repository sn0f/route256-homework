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

func TestPurchase(t *testing.T) {
	type in struct {
		ctx  context.Context
		user int64
	}

	type out struct {
		orderID int64
		err     error
	}

	var (
		ctx         = context.Background()
		workerCount = 5

		sqlError                = errors.New("sql error")
		serviceUnavailableError = errors.New("service unavailable")
		listCartRepoError       = fmt.Errorf("listing cart: %v", sqlError)
		listCartError           = fmt.Errorf("list cart: %v", listCartRepoError)
		createOrderError        = fmt.Errorf("creating order: %v", serviceUnavailableError)
		clearCartError          = fmt.Errorf("clearing cart: %v", sqlError)

		product1 = domain.Product{
			Name:  gofakeit.BeerName(),
			Price: 10,
		}

		sku1 = gofakeit.Uint32()

		cartItems = []domain.CartItem{
			{
				SKU:   sku1,
				Count: 5,
				Name:  product1.Name,
				Price: product1.Price,
			},
		}

		emptyOrder = domain.Order{
			OrderID: int64(0),
		}
		order = domain.Order{
			OrderID: gofakeit.Int64(),
		}
	)

	tests := []struct {
		name   string
		in     in
		setup  func(m *mocksService, in *in)
		assert func(require.TestingT, *in, *out)
	}{
		{
			name: "negative case - listing cart error",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return([]domain.CartItem{}, sqlError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.orderID, emptyOrder.OrderID)
				require.Equal(t, out.err, listCartError)
			},
		},
		{
			name: "negative case - create order error",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return(cartItems, nil)
				m.products.GetProductMock.When(ctx, sku1).Then(&product1, nil)
				m.orders.CreateOrderMock.Expect(ctx, in.user, cartItems).Return(&emptyOrder, serviceUnavailableError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.orderID, emptyOrder.OrderID)
				require.Equal(t, out.err, createOrderError)
			},
		},
		{
			name: "negative case - clear cart error",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return(cartItems, nil)
				m.repo.ClearCartMock.Expect(ctx, in.user).Return(sqlError)
				m.products.GetProductMock.When(ctx, sku1).Then(&product1, nil)
				m.orders.CreateOrderMock.Expect(ctx, in.user, cartItems).Return(&order, nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, out.orderID, emptyOrder.OrderID)
				require.Equal(t, out.err, clearCartError)
			},
		},
		{
			name: "positive case",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return(cartItems, nil)
				m.repo.ClearCartMock.Expect(ctx, in.user).Return(nil)
				m.products.GetProductMock.When(ctx, sku1).Then(&product1, nil)
				m.orders.CreateOrderMock.Expect(ctx, in.user, cartItems).Return(&order, nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, order.OrderID, out.orderID)
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

			cart, err := model.Purchase(ctx, tt.in.user)
			tt.assert(t, &tt.in, &out{
				orderID: cart,
				err:     err,
			})
		})
	}
}
