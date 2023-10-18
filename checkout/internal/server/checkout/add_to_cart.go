package checkout

import (
	"context"
	"errors"
	"fmt"
	api "route256/checkout/pkg/checkout/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func ValidateAddToCartRequest(r *api.AddToCartRequest) error {
	var errs []error

	if r.GetUser() <= 0 {
		errs = append(errs, ErrEmptyUser)
	}
	if r.GetSku() == 0 {
		errs = append(errs, ErrEmptySKU)
	}
	if r.GetCount() == 0 {
		errs = append(errs, ErrEmptyCount)
	}

	return errors.Join(errs...)
}

func (s *checkoutServer) AddToCart(ctx context.Context, req *api.AddToCartRequest) (*emptypb.Empty, error) {
	err := ValidateAddToCartRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validating request: %v", err)
	}

	err = s.model.AddToCart(ctx, req.User, req.Sku, uint16(req.Count))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
