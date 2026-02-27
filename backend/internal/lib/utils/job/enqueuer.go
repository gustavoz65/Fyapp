package job

import "github.com/hibiken/asynq"

// Enqueuer define a interface para enfileirar tarefas de job
// Esta interface evita ciclos de importação entre service e job packages
type Enqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}
