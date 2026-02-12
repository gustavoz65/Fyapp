package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TaskProcessRecurrings = "recurring:process"
)

type ProcessRecurringsPayload struct{}

func NewProcessRecurringsTask() (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessRecurringsPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskProcessRecurrings, payload, asynq.Queue("default")), nil
}

func (j *JobService) handleProcessRecurringsTask(ctx context.Context, t *asynq.Task) error {
	j.logger.Info().Msg("Processing due recurring transactions")

	if j.recurringService == nil {
		return fmt.Errorf("recurring service not initialized")
	}

	err := j.recurringService.ProcessDueRecurrings(ctx)
	if err != nil {
		j.logger.Error().Err(err).Msg("Failed to process recurring transactions")
		return err
	}

	j.logger.Info().Msg("Successfully processed recurring transactions")
	return nil
}
