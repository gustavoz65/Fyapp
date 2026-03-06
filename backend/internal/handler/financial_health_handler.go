package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/service"
)

type FinancialHealthHandler struct {
	healthService *service.FinancialHealthService
}

func NewFinancialHealthHandler(healthService *service.FinancialHealthService) *FinancialHealthHandler {
	return &FinancialHealthHandler{
		healthService: healthService,
	}
}

// GetCurrentScore retorna o score do mês atual
func (h *FinancialHealthHandler) GetCurrentScore(c echo.Context) error {
	userID := middleware.GetUserID(c)

	snapshot, err := h.healthService.GetSnapshot(c.Request().Context(), userID, time.Now())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, snapshot)
}

// GetHistoricalScores retorna histórico de scores
func (h *FinancialHealthHandler) GetHistoricalScores(c echo.Context) error {
	userID := middleware.GetUserID(c)

	snapshots, err := h.healthService.GetHistoricalSnapshots(c.Request().Context(), userID, 6)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"snapshots": snapshots,
	})
}
