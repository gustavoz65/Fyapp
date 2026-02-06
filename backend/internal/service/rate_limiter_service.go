package service

import (
	"github.com/gustavoz65/Cashing-go/internal/repository"
	"golang.org/x/time/rate"
)

type RateLimiterService struct {
	defaultLimiter repository.RateLimiter
}

func NewRateLimiterService(requestsPerSecond float64, burst int) *RateLimiterService {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	wrapper := repository.NewLimiterWrapper(limiter)

	return &RateLimiterService{
		defaultLimiter: wrapper,
	}
}

func (s *RateLimiterService) GetDefaultLimiter() repository.RateLimiter {
	return s.defaultLimiter
}

func (s *RateLimiterService) CreateLimiter(requestsPerSecond float64, burst int) repository.RateLimiter {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	return repository.NewLimiterWrapper(limiter)
}
