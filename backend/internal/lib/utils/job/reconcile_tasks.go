package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TaskAutoReconcile = "reconcile:auto"
)

type AutoReconcilePayload struct{}

func NewAutoReconcileTask() (*asynq.Task, error) {
	payload, err := json.Marshal(AutoReconcilePayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskAutoReconcile, payload, asynq.Queue("default")), nil
}

func (j *JobService) handleAutoReconcileTask(ctx context.Context, t *asynq.Task) error {
	j.logger.Info().Msg("Running auto-reconciliation for due transactions")

	if j.transactionService == nil {
		return fmt.Errorf("transaction service not initialized")
	}

	if err := j.transactionService.AutoReconcile(ctx); err != nil {
		j.logger.Error().Err(err).Msg("Failed to auto-reconcile transactions")
		return err
	}

	j.logger.Info().Msg("Auto-reconciliation completed")
	return nil
}
