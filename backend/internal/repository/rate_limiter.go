package repository

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	Max         int
	Duration    time.Duration
	KeyFunc     func(r *fiber.Ctx) string
	Endpoint    string
	SkipOnerror bool
}
