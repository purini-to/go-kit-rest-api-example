package services

import (
	"context"
	"github.com/purini-to/go-kit-rest-api-example/entities"
)

type Service interface {
	GetTasks(ctx context.Context) ([]entities.Task, error)
	GetTask(_ context.Context, id string) (entities.Task, error)
}

type service struct {
	*task
}

func NewService() Service {
	return &service{
		task: newTask(),
	}
}
