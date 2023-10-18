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

func TestListCart(t *testing.T) {
	type in struct {
		ctx  context.Context
		user int64
	}

	type out struct {
		cart domain.Cart
		err  error
	}

	var (
		ctx         = context.Background()
		workerCount = 5

		sqlError                = errors.New("sql error")
		serviceUnavailableError = errors.New("service unavailable")
		getProductError         = fmt.Errorf("get product: %v", serviceUnavailableError)
		listCartRepoError       = fmt.Errorf("listing cart: %v", sqlError)
		listCartRepoError2      = fmt.Errorf("listing cart: %v", getProductError)

		product1 = domain.Product{
			Name:  gofakeit.BeerName(),
			Price: 10,
		}
		product2 = domain.Product{
			Name:  gofakeit.BeerName(),
			Price: 20,
		}

		sku1 = gofakeit.Uint32()
		sku2 = gofakeit.Uint32()

		cartItems = []domain.CartItem{
			{
				SKU:   sku1,
				Count: 5,
				Name:  product1.Name,
				Price: product1.Price,
			},
			{
				SKU:   sku2,
				Count: 3,
				Name:  product2.Name,
				Price: product2.Price,
			},
		}

		emptyCart     = domain.Cart{}
		cartWithItems = domain.Cart{
			Items:      cartItems,
			TotalPrice: 110,
		}
	)

	tests := []struct {
		name   string
		in     in
		setup  func(m *mocksService, in *in)
		assert func(require.TestingT, *in, *out)
	}{
		{
			name: "negative case - listing cart repo error",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return([]domain.CartItem{}, sqlError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, len(emptyCart.Items), len(out.cart.Items))
				require.Equal(t, emptyCart.TotalPrice, out.cart.TotalPrice)
				require.Equal(t, out.err, listCartRepoError)
			},
		},
		{
			name: "negative case - get product error",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return(cartItems, nil)
				m.products.GetProductMock.When(ctx, sku1).Then(&domain.Product{}, serviceUnavailableError)
				m.products.GetProductMock.When(ctx, sku2).Then(&domain.Product{}, serviceUnavailableError)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, len(emptyCart.Items), len(out.cart.Items))
				require.Equal(t, out.err, listCartRepoError2)
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
				m.products.GetProductMock.When(ctx, sku1).Then(&product1, nil)
				m.products.GetProductMock.When(ctx, sku2).Then(&product2, nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, len(cartWithItems.Items), len(out.cart.Items))
				require.Equal(t, cartWithItems.TotalPrice, out.cart.TotalPrice)
				require.NoError(t, out.err)
			},
		},
		{
			name: "positive case - empty cart",
			in: in{
				ctx:  ctx,
				user: gofakeit.Int64(),
			},
			setup: func(m *mocksService, in *in) {
				m.repo.ListCartMock.Expect(ctx, in.user).Return([]domain.CartItem{}, nil)
			},
			assert: func(t require.TestingT, in *in, out *out) {
				require.Equal(t, len(emptyCart.Items), len(out.cart.Items))
				require.Equal(t, emptyCart.TotalPrice, out.cart.TotalPrice)
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

			cart, err := model.ListCart(ctx, tt.in.user)
			tt.assert(t, &tt.in, &out{
				cart: cart,
				err:  err,
			})
		})
	}
}
