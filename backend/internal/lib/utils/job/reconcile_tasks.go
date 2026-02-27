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
	j.Logger.Info().Msg("Running auto-reconciliation for due transactions")

	if j.TransactionService == nil {
		return fmt.Errorf("transaction service not initialized")
	}

	if err := j.TransactionService.AutoReconcile(ctx); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to auto-reconcile transactions")
		return err
	}

	j.Logger.Info().Msg("Auto-reconciliation completed")
	return nil
}
