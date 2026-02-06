package middleware

import (
	"net/http"

	"github.com/gustavoz65/Cashing-go/internal/repository"
)

func RateLimitMiddleware(next http.Handler, rl repository.RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
