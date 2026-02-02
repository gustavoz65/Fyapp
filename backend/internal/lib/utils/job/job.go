package job

import (
	"github.com/gustavoz65/Cashing-go/internal/config"
	"github.com/hibiken/asynq"
	zerolog "github.com/rs/zerolog"
)

type JobService struct {
	Client *asynq.Client
	Server *asynq.Server
	logger *zerolog.Logger
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
				"critical": 6, // Maior Prioridade
				"default":  3, // Valor Default
				"low":      1, // Menor Prioridade
			},
		},
	)
	return &JobService{
		Client: client,
		Server: server,
		logger: logger,
	}
}

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.logger.Info().Msg("Starting background job server")
	if err := j.Server.Start(mux); err != nil {
		return err
	}
	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.Server.Shutdown()
	j.Client.Close()
}
