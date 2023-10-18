package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger, serviceName string) {
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}

	_, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}
