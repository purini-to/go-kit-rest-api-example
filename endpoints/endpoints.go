package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/purini-to/go-kit-rest-api-example/entities"
	"github.com/purini-to/go-kit-rest-api-example/services"
)

func GetTasksEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		t, e := s.GetTasks(ctx)
		return GetTasksOutput{Tasks: t}, e
	}
}

func GetTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTaskInput)
		t, e := s.GetTask(ctx, req.ID)
		return GetTaskOutput{Task: t}, e
	}
}

func PostTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostTaskInput)
		t, e := s.PostTask(ctx, req.Task)
		return PostTaskOutput{Task: t}, e
	}
}

func PutTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutTaskInput)
		t, e := s.PutTask(ctx, req.ID, req.Task)
		return PutTaskOutput{Task: t}, e
	}
}

func DeleteTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTaskInput)
		e := s.DeleteTask(ctx, req.ID)
		return nil, e
	}
}

type GetTasksOutput struct {
	Tasks []entities.Task
}

type GetTaskInput struct {
	ID string
}

type GetTaskOutput struct {
	Task entities.Task
}

type PostTaskInput struct {
	Task entities.Task
}

type PostTaskOutput struct {
	Task entities.Task
}

type PutTaskInput struct {
	ID   string
	Task entities.Task
}

type PutTaskOutput struct {
	Task entities.Task
}

type DeleteTaskInput struct {
	ID   string
	Task entities.Task
}
