package repository

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
}

type LimiterWrapper struct {
	limiter *rate.Limiter
}

func NewLimiterWrapper(limiter *rate.Limiter) *LimiterWrapper {
	return &LimiterWrapper{limiter: limiter}
}

func (l *LimiterWrapper) Allow() bool {
	return l.limiter.Allow()
}

func (l *LimiterWrapper) Wait(ctx context.Context) error {
	return l.limiter.Wait(ctx)
}
