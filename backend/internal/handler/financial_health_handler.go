package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/repository"
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
	ctx := c.Request().Context()
	now := time.Now()

	// Tenta buscar snapshot existente
	snapshot, err := h.healthService.GetSnapshot(ctx, userID, now)
	if err != nil {
		// Se não encontrou, calcula um novo
		if errors.Is(err, repository.ErrSnapshotNotFound) {
			snapshot, err = h.healthService.CalculateScore(ctx, userID, now)
			if err != nil {
				return err
			}
		} else {
			return err
		}
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
