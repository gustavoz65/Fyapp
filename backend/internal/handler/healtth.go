package handler

import (
	"context"
	"time"

	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/your_project/backend/internal/server"
)

type HealthHandler struct {
	Handler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		Handler: NewHandler(s),
	}
}

func (h *HealthHandler) CheckHandler(c *echo.Context) error {
	start := time.Now()
	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	response := map[string]interface{}{
		"status":      "ok",
		"timestamp":   time.Now().UTC(),
		"environment": h.server.Config.Primary.Env,
		"checks":      make(map[string]interface{}),
	}

	checks := response["checks"].(map[string]interface{})
	isHealthy := true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	if err != h.server.DB.Pool.Ping(ctx); err != nil {
		checks["database"] = map[string]interface{}{
			"status":       "error",
			"reponse_time": time.Since(dbStart).String(),
			"error":        err.Error(),
		}
		isHealthy = false
		logger.Error().Err(err).Dur("reponse_time", time.Since(dbStart).Msg("o Banco de Dados não está Saudavel"))
		if h.server.LoggerService != nil && h.server.LoggerService.GetApplication() != nil {
			h.server.LoggerService.GetApplication().RecordCustomEvent(
				"HealthCheckError", map[string]interface{}{
					"check_type": "database",
					"operation": "health_check",
					"error_type": "database_unhealthy",
					"error_message": err.Error(),
				}
			)
		}
	}
}
