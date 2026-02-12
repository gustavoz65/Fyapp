package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/repository"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, config repository.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if rdb == nil {
			return c.Next()
		}
		ctx := c.Context()

		key := ""
		if config.KeyFunc != nil {
			key = config.KeyFunc(c)
		}
		if key == "" {
			key = fmt.Sprintf("ip:%s:path:%s", c.IP(), c.Path())
		}
		redisKey := fmt.Sprintf("ratelimit:%s", key)

		count, err := rdb.Get(ctx, redisKey).Int()
		if err != nil && err != redis.Nil {
			if config.SkipOnerror {
				return c.Next()
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error checking rate limit",
			})
		}

		if count >= config.Max {
			c.Set("Retry-After", fmt.Sprintf("%.0f", config.Duration.Seconds()))
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
			c.Set("X-RateLimit-Remaining", "0")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
			})
		}

		pipe := rdb.Pipeline()
		pipe.Incr(ctx, redisKey)
		pipe.Expire(ctx, redisKey, config.Duration)
		_, _ = pipe.Exec(ctx)

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Max))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", config.Max-count-1))

		return c.Next()
	}
}

func DefaultRateLimit(rdb *redis.Client) fiber.Handler {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      300,
		Duration: time.Minute,
		KeyFunc: func(c *fiber.Ctx) string {
			return fmt.Sprintf("ip:%s:path:%s", c.IP(), c.Path())
		},
		Endpoint:    "global",
		SkipOnerror: true,
	})
}

func UserRateLimit(rdb *redis.Client, authCheckFunc func(*fiber.Ctx) error) fiber.Handler {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      1000,
		Duration: time.Minute,
		KeyFunc: func(c *fiber.Ctx) string {
			if authCheckFunc != nil && authCheckFunc(c) == nil {
				if userID, ok := c.Locals("userID").(uuid.UUID); ok {
					return fmt.Sprintf("user:%s", userID.String())
				}
			}
			return fmt.Sprintf("ip:%s", c.IP())
		},
		Endpoint:    "user",
		SkipOnerror: true,
	})
}

func AuthRateLimit(rdb *redis.Client) fiber.Handler {
	return RateLimit(rdb, repository.RateLimiter{
		Max:      5,
		Duration: 15 * time.Minute,
		KeyFunc: func(c *fiber.Ctx) string {
			return fmt.Sprintf("auth:%s", c.IP())
		},
		Endpoint:    "auth",
		SkipOnerror: true,
	})
}
