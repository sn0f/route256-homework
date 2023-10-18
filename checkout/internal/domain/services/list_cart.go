package services

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/libs/server/wpool"
)

func (m *Model) ListCart(ctx context.Context, user int64) (domain.Cart, error) {
	items, err := m.cartRepository.ListCart(ctx, user)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("listing cart: %v", err)
	}

	tasks := make([]wpool.Task[domain.CartItem, domain.CartItem], 0, len(items))

	for _, cartItem := range items {
		fn := func(ctx context.Context, item domain.CartItem) (domain.CartItem, error) {
			product, err := m.productManager.GetProduct(ctx, item.SKU)
			if err != nil {
				return item, fmt.Errorf("get product: %v", err)
			}

			item.Name = product.Name
			item.Price = product.Price
			return item, nil
		}

		task := wpool.Task[domain.CartItem, domain.CartItem]{
			Func: fn,
			Args: cartItem,
		}
		tasks = append(tasks, task)
	}

	cart := domain.Cart{
		Items:      make([]domain.CartItem, 0, len(items)),
		TotalPrice: 0,
	}

	wp := wpool.NewWorkerPool[domain.CartItem, domain.CartItem](m.listCartWorkerCount)

	go wp.StartTasks(ctx, tasks)

	for r := range wp.Results() {
		if r.Error != nil {
			return domain.Cart{}, fmt.Errorf("listing cart: %v", r.Error.Error())
		}

		item := r.Data
		cart.Items = append(cart.Items, item)
		cart.TotalPrice += item.Price * uint32(item.Count)
	}

	return cart, nil
}
