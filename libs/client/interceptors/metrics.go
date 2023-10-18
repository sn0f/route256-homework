package interceptors

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "client",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"handler", "status"},
	)
)

func MetricsInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	timeStart := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	code := status.Convert(err).Code().String()

	if err != nil {
		return err
	}

	elapsed := time.Since(timeStart)
	HistogramResponseTime.WithLabelValues(method, code).Observe(elapsed.Seconds())
	return nil
}
