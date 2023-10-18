package product

import (
	"context"
	"route256/checkout/internal/domain"
	api "route256/checkout/pkg/product/v1"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	histogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "get_product",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"source"},
	)
)

func (c *client) GetProduct(ctx context.Context, sku uint32) (*domain.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetProduct processing")
	defer span.Finish()

	span.SetTag("sku", sku)

	timeStart := time.Now()
	product, ok := c.getProductCache.Get(ctx, sku)
	if ok {
		elapsed := time.Since(timeStart)
		histogramResponseTime.WithLabelValues("cache").Observe(elapsed.Seconds())
		return &product, nil
	}

	timeStart = time.Now()

	request := &api.GetProductRequest{
		Sku:   sku,
		Token: c.Token,
	}

	c.Limiter.Wait(ctx)

	response, err := c.productClient.GetProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	product = domain.Product{
		Name:  response.Name,
		Price: response.Price,
	}

	elapsed := time.Since(timeStart)
	histogramResponseTime.WithLabelValues("grpc").Observe(elapsed.Seconds())

	c.getProductCache.Set(ctx, sku, product)

	return &product, nil
}
