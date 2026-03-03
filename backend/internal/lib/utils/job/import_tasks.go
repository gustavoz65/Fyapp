package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TaskImportTransactions = "import:transactions"
)

type ImportTransactionsPayload struct {
	JobID         string `json:"job_id"`
	UserID        string `json:"user_id"`
	BankAccountID string `json:"bank_account_id"`
	BankType      string `json:"bank_type"`
	ForceReimport bool   `json:"force_reimport"`
	CSVData       string `json:"csv_data"`
}

type ImportStatus struct {
	JobID      string   `json:"job_id"`
	Status     string   `json:"status"` // "processing", "completed", "failed"
	Progress   int      `json:"progress"`
	Total      int      `json:"total"`
	Imported   int      `json:"imported"`
	Duplicates int      `json:"duplicates"`
	Errors     []string `json:"errors"`
	Message    string   `json:"message"`
}

// NewImportTransactionsTask creates a new import task
func NewImportTransactionsTask(payload ImportTransactionsPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal import payload: %w", err)
	}
	return asynq.NewTask(TaskImportTransactions, data), nil
}

// handleImportTransactionsTask processes the import job
func (j *JobService) handleImportTransactionsTask(ctx context.Context, t *asynq.Task) error {
	var p ImportTransactionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal import payload: %w", err)
	}

	j.Logger.Info().
		Str("job_id", p.JobID).
		Str("user_id", p.UserID).
		Msg("Starting transaction import job")

	// Update status to processing
	j.updateImportStatus(p.JobID, ImportStatus{
		JobID:   p.JobID,
		Status:  "processing",
		Message: "Processando arquivo CSV...",
	})

	// Parse UUIDs
	userID, err := uuid.Parse(p.UserID)
	if err != nil {
		j.updateImportStatus(p.JobID, ImportStatus{
			JobID:   p.JobID,
			Status:  "failed",
			Message: fmt.Sprintf("Invalid user ID: %v", err),
		})
		return fmt.Errorf("invalid user ID: %w", err)
	}

	bankAccountID, err := uuid.Parse(p.BankAccountID)
	if err != nil {
		j.updateImportStatus(p.JobID, ImportStatus{
			JobID:   p.JobID,
			Status:  "failed",
			Message: fmt.Sprintf("Invalid bank account ID: %v", err),
		})
		return fmt.Errorf("invalid bank account ID: %w", err)
	}

	// Create reader from CSV data
	reader := strings.NewReader(p.CSVData)

	// Call the import service with progress callback
	result, err := j.TransactionService.ImportTransactionsWithProgress(
		ctx,
		userID,
		bankAccountID,
		reader,
		p.BankType,
		p.ForceReimport,
		func(current, total int) {
			// Update progress in real-time
			progress := 0
			if total > 0 {
				progress = (current * 100) / total
			}
			j.updateImportStatus(p.JobID, ImportStatus{
				JobID:    p.JobID,
				Status:   "processing",
				Progress: progress,
				Total:    total,
				Message:  fmt.Sprintf("Processando transação %d de %d...", current, total),
			})
		},
	)

	if err != nil {
		j.Logger.Error().Err(err).
			Str("job_id", p.JobID).
			Msg("Failed to import transactions")

		j.updateImportStatus(p.JobID, ImportStatus{
			JobID:   p.JobID,
			Status:  "failed",
			Message: fmt.Sprintf("Erro ao importar: %v", err),
		})
		return err
	}

	// Success
	j.updateImportStatus(p.JobID, ImportStatus{
		JobID:      p.JobID,
		Status:     "completed",
		Progress:   100,
		Total:      result.TotalImported + result.Duplicates,
		Imported:   result.TotalImported,
		Duplicates: result.Duplicates,
		Errors:     result.Errors,
		Message:    fmt.Sprintf("Importação concluída! %d novas, %d duplicadas", result.TotalImported, result.Duplicates),
	})

	j.Logger.Info().
		Str("job_id", p.JobID).
		Int("imported", result.TotalImported).
		Int("duplicates", result.Duplicates).
		Msg("Import job completed successfully")

	return nil
}

// updateImportStatus updates the import status in Redis
func (j *JobService) updateImportStatus(jobID string, status ImportStatus) {
	ctx := context.Background()

	if err := j.SetImportStatus(ctx, jobID, status); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to set import status")
		return
	}

	j.Logger.Debug().
		Str("job_id", jobID).
		Str("status", status.Status).
		Int("progress", status.Progress).
		Msg("Import status updated")
}
