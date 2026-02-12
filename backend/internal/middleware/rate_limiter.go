package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, config repository.RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rdb == nil {
				return next(c)
			}
			ctx := c.Request().Context()

			key := ""
			if config.KeyFunc != nil {
				key = config.KeyFunc(c)
			}
			if key == "" {
				key = fmt.Sprintf("ip:%s:path:%s", c.RealIP(), c.Path())
			}
			redisKey := fmt.Sprintf("ratelimit:%s", key)

			count, err := rdb.Get(ctx, redisKey).Int()
			if err != nil && err != redis.Nil {
				if config.SkipOnerror {
					return next(c)
				}
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Error checking rate limit",
				})
			}

			if count >= config.Max {
				c.Response().Header().Set("Retry-After", fmt.Sprintf("%.0f", config.Duration.Seconds()))
				c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")

				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"success": false,
					"error":   "Rate limit exceeded. Please try again later.",
				})
			}

			pipe := rdb.Pipeline()
			pipe.Incr(ctx, redisKey)
			pipe.Expire(ctx, redisKey, config.Duration)
			_, _ = pipe.Exec(ctx)

			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", config.Max-count-1))

			return next(c)
		}
	}
}

func DefaultRateLimit(rdb *redis.Client) echo.MiddlewareFunc {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      300,
		Duration: time.Minute,
		KeyFunc: func(c echo.Context) string {
			return fmt.Sprintf("ip:%s:path:%s", c.RealIP(), c.Path())
		},
		Endpoint:    "global",
		SkipOnerror: true,
	})
}

func UserRateLimit(rdb *redis.Client) echo.MiddlewareFunc {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      1000,
		Duration: time.Minute,
		KeyFunc: func(c echo.Context) string {
			if userID := c.Get("user_id"); userID != nil {
				if uid, ok := userID.(uuid.UUID); ok {
					return fmt.Sprintf("user:%s", uid.String())
				}
			}
			return fmt.Sprintf("ip:%s", c.RealIP())
		},
		Endpoint:    "user",
		SkipOnerror: true,
	})
}

func AuthRateLimit(rdb *redis.Client) echo.MiddlewareFunc {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      5,
		Duration: 15 * time.Minute,
		KeyFunc: func(c echo.Context) string {
			return fmt.Sprintf("auth:%s", c.RealIP())
		},
		Endpoint:    "auth",
		SkipOnerror: true,
	})
}

// MutationRateLimit aplica rate limiting para operações de escrita (POST/PUT/PATCH/DELETE)
func MutationRateLimit(rdb *redis.Client) echo.MiddlewareFunc {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      100, // 100 mutations por minuto
		Duration: time.Minute,
		KeyFunc: func(c echo.Context) string {
			// Usar user_id se autenticado, senão IP
			if userID := c.Get("user_id"); userID != nil {
				if uid, ok := userID.(uuid.UUID); ok {
					return fmt.Sprintf("mutations:user:%s", uid.String())
				}
			}
			return fmt.Sprintf("mutations:ip:%s", c.RealIP())
		},
		Endpoint:    "mutations",
		SkipOnerror: true,
	})
}

// ReadRateLimit aplica rate limiting para operações de leitura (GET)
func ReadRateLimit(rdb *redis.Client) echo.MiddlewareFunc {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      500, // 500 reads por minuto
		Duration: time.Minute,
		KeyFunc: func(c echo.Context) string {
			// Usar user_id se autenticado, senão IP
			if userID := c.Get("user_id"); userID != nil {
				if uid, ok := userID.(uuid.UUID); ok {
					return fmt.Sprintf("reads:user:%s", uid.String())
				}
			}
			return fmt.Sprintf("reads:ip:%s", c.RealIP())
		},
		Endpoint:    "reads",
		SkipOnerror: true,
	})
}
