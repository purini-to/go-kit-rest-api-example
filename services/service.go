package services

import (
	"context"
	"github.com/purini-to/go-kit-rest-api-example/entities"
)

type Service interface {
	GetTasks(ctx context.Context) ([]entities.Task, error)
	GetTask(_ context.Context, id string) (entities.Task, error)
	PostTask(ctx context.Context, t entities.Task) (entities.Task, error)
	PutTask(ctx context.Context, id string, t entities.Task) (entities.Task, error)
	DeleteTask(ctx context.Context, id string) error
	HealthCheck(_ context.Context) string
}

type service struct {
	*task
	*health
}

func NewService() Service {
	return &service{
		task: newTask(),
		health: newHealth(),
	}
}
