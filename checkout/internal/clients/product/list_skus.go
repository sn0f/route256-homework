package product

import (
	"context"
	api "route256/checkout/pkg/product/v1"

	"github.com/opentracing/opentracing-go"
)

func (c *client) ListSKUs(ctx context.Context, startAfterSku uint32, count uint32) ([]uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ListSKUs processing")
	defer span.Finish()

	span.SetTag("startAfterSku", startAfterSku)
	span.SetTag("count", count)

	request := &api.ListSkusRequest{
		StartAfterSku: startAfterSku,
		Count:         count,
		Token:         c.Token,
	}

	response, err := c.productClient.ListSkus(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Skus, nil
}
