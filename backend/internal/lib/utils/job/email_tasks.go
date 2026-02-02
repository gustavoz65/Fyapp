package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskWelcome = "email:welcome"
)

type WelcomeEmailPlayload struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
}

func NewWelcomeEmailTasl(to, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPlayload{
		To:        to,
		FirstName: firstName,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second)), nil
}
