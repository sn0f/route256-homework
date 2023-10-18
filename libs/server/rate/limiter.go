package rate

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) (err error)
}

func NewLimiter(rpsLimit uint64, burstLimit uint64) RateLimiter {
	return rate.NewLimiter(rate.Limit((rpsLimit)), int(burstLimit))
}
