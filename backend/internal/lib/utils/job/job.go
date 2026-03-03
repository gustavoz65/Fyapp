package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	zerolog "github.com/rs/zerolog"

	"github.com/gustavoz65/Fyapp/internal/config"
	"github.com/gustavoz65/Fyapp/internal/service"
)

type JobService struct {
	Client             *asynq.Client
	Server             *asynq.Server
	Scheduler          *asynq.Scheduler
	Redis              *redis.Client
	Logger             *zerolog.Logger
	RecurringService   *service.RecurringTransactionService
	TransactionService *service.TransactionService
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	location := time.UTC
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{Location: location},
	)

	// Create Redis client for status storage
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	return &JobService{
		Client:    client,
		Server:    server,
		Scheduler: scheduler,
		Redis:     redisClient,
		Logger:    logger,
	}
}

func (j *JobService) SetRecurringService(svc *service.RecurringTransactionService) {
	j.RecurringService = svc
}

func (j *JobService) SetTransactionService(svc *service.TransactionService) {
	j.TransactionService = svc
}

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)
	mux.HandleFunc(TaskProcessRecurrings, j.handleProcessRecurringsTask)
	mux.HandleFunc(TaskAutoReconcile, j.handleAutoReconcileTask)
	mux.HandleFunc(TaskImportTransactions, j.handleImportTransactionsTask)

	task, _ := NewProcessRecurringsTask()
	if _, err := j.Scheduler.Register("0 0 * * *", task); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to schedule recurring task")
	}

	reconcileTask, _ := NewAutoReconcileTask()
	if _, err := j.Scheduler.Register("5 0 * * *", reconcileTask); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to schedule auto reconcile task")
	}

	go func() {
		if err := j.Scheduler.Run(); err != nil {
			j.Logger.Error().Err(err).Msg("Scheduler stopped")
		}
	}()

	j.Logger.Info().Msg("Starting background job server")
	if err := j.Server.Start(mux); err != nil {
		return err
	}
	return nil
}

func (j *JobService) Stop() {
	j.Logger.Info().Msg("Stopping background job server")
	j.Scheduler.Shutdown()
	j.Server.Shutdown()
	if err := j.Client.Close(); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to close Asynq client")
	}
	if err := j.Redis.Close(); err != nil {
		j.Logger.Error().Err(err).Msg("Failed to close Redis client")
	}
}

// SetImportStatus stores import status in Redis
func (j *JobService) SetImportStatus(ctx context.Context, jobID string, status interface{}) error {
	key := fmt.Sprintf("import_status:%s", jobID)
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	// Store with 1 hour expiration
	return j.Redis.Set(ctx, key, data, time.Hour).Err()
}

// GetImportStatus retrieves import status from Redis
func (j *JobService) GetImportStatus(ctx context.Context, jobID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("import_status:%s", jobID)
	data, err := j.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var status map[string]interface{}
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return status, nil
}
