package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"

	"github.com/gustavoz65/Fyapp/internal/lib/utils/job"
	"github.com/gustavoz65/Fyapp/internal/middleware"
)

type WebSocketHandler struct {
	jobService *job.JobService
	logger     *zerolog.Logger
}

func NewWebSocketHandler(jobService *job.JobService, logger *zerolog.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		jobService: jobService,
		logger:     logger,
	}
}

// ImportProgress handles WebSocket connections for import progress updates
func (h *WebSocketHandler) ImportProgress(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			if err := ws.Close(); err != nil {
				h.logger.Error().Err(err).Msg("Failed to close WebSocket connection")
			}
		}()

		// Get job ID from query param
		jobID := c.QueryParam("job_id")
		if jobID == "" {
			h.logger.Error().Msg("WebSocket: job_id is required")
			return
		}

		// Get user ID from context (already authenticated via middleware)
		userID := middleware.GetUserID(c)
		h.logger.Info().
			Str("user_id", userID.String()).
			Str("job_id", jobID).
			Msg("WebSocket connection established for import progress")

		ctx := context.Background()
		ticker := time.NewTicker(500 * time.Millisecond) // Poll every 500ms
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Get import status from Redis
				status, err := h.jobService.GetImportStatus(ctx, jobID)
				if err != nil {
					// Job not found or expired
					errMsg := map[string]interface{}{
						"error":  "Job não encontrado ou expirado",
						"job_id": jobID,
						"status": "error",
					}
					data, _ := json.Marshal(errMsg)
					if err := websocket.Message.Send(ws, string(data)); err != nil {
						h.logger.Error().Err(err).Msg("Failed to send error message")
					}
					return
				}

				// Send status update
				data, err := json.Marshal(status)
				if err != nil {
					h.logger.Error().Err(err).Msg("Failed to marshal status")
					continue
				}

				if err := websocket.Message.Send(ws, string(data)); err != nil {
					h.logger.Error().Err(err).Msg("Failed to send status update")
					return
				}

				// Check if import is completed or failed
				if statusStr, ok := status["status"].(string); ok {
					if statusStr == "completed" || statusStr == "failed" {
						h.logger.Info().
							Str("job_id", jobID).
							Str("status", statusStr).
							Msg("Import job finished, closing WebSocket")
						time.Sleep(1 * time.Second) // Give client time to receive final status
						return
					}
				}

			case <-ws.Request().Context().Done():
				h.logger.Info().
					Str("job_id", jobID).
					Msg("WebSocket connection closed by client")
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
