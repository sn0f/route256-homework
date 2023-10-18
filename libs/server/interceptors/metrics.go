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
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "grpc",
		Name:      "requests_total",
	},
		[]string{"handler"},
	)
	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "grpc",
		Name:      "responses_total",
	},
		[]string{"handler", "status"},
	)
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"handler", "status"},
	)
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	RequestsCounter.WithLabelValues(info.FullMethod).Inc()

	timeStart := time.Now()
	res, err := handler(ctx, req)
	code := status.Convert(err).Code().String()

	if err != nil {
		ResponseCounter.WithLabelValues(info.FullMethod, code).Inc()
		return nil, err
	}

	elapsed := time.Since(timeStart)
	HistogramResponseTime.WithLabelValues(info.FullMethod, code).Observe(elapsed.Seconds())
	ResponseCounter.WithLabelValues(info.FullMethod, code).Inc()
	return res, nil
}
