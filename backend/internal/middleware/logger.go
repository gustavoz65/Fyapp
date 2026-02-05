package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const loggerKey = "logger"

// LoggerMiddleware injeta o logger no contexto e loga requisicoes
func LoggerMiddleware(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()

			// Injeta logger no contexto
			requestLogger := logger.With().
				Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("remote_ip", c.RealIP()).
				Logger()

			c.Set(loggerKey, &requestLogger)

			err := next(c)

			duration := time.Since(start)
			status := c.Response().Status

			logEvent := requestLogger.Info()
			if status >= 500 {
				logEvent = requestLogger.Error()
			} else if status >= 400 {
				logEvent = requestLogger.Warn()
			}

			logEvent.
				Int("status", status).
				Dur("duration", duration).
				Int64("bytes_out", c.Response().Size).
				Msg("request completed")

			return err
		}
	}
}

// GetLogger retorna o logger do contexto
func GetLogger(c echo.Context) *zerolog.Logger {
	if logger, ok := c.Get(loggerKey).(*zerolog.Logger); ok {
		return logger
	}
	l := zerolog.Nop()
	return &l
}
