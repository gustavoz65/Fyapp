package job

import (
	"time"

	"github.com/gustavoz65/Cashing-go/internal/config"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/hibiken/asynq"
	zerolog "github.com/rs/zerolog"
)

type JobService struct {
	Client             *asynq.Client
	Server             *asynq.Server
	Scheduler          *asynq.Scheduler
	logger             *zerolog.Logger
	recurringService   *service.RecurringTransactionService
	transactionService *service.TransactionService
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
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
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{Location: location},
	)

	return &JobService{
		Client:    client,
		Server:    server,
		Scheduler: scheduler,
		logger:    logger,
	}
}

func (j *JobService) SetRecurringService(svc *service.RecurringTransactionService) {
	j.recurringService = svc
}

func (j *JobService) SetTransactionService(svc *service.TransactionService) {
	j.transactionService = svc
}

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)
	mux.HandleFunc(TaskProcessRecurrings, j.handleProcessRecurringsTask)
	mux.HandleFunc(TaskAutoReconcile, j.handleAutoReconcileTask)

	task, _ := NewProcessRecurringsTask()
	if _, err := j.Scheduler.Register("0 0 * * *", task); err != nil {
		j.logger.Error().Err(err).Msg("Failed to schedule recurring task")
	}

	reconcileTask, _ := NewAutoReconcileTask()
	if _, err := j.Scheduler.Register("5 0 * * *", reconcileTask); err != nil {
		j.logger.Error().Err(err).Msg("Failed to schedule auto reconcile task")
	}

	go func() {
		if err := j.Scheduler.Run(); err != nil {
			j.logger.Error().Err(err).Msg("Scheduler stopped")
		}
	}()

	j.logger.Info().Msg("Starting background job server")
	if err := j.Server.Start(mux); err != nil {
		return err
	}
	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.Scheduler.Shutdown()
	j.Server.Shutdown()
	j.Client.Close()
}
