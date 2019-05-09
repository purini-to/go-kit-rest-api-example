package go_kit_rest_api_example

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/purini-to/go-kit-rest-api-example/entities"
	"github.com/purini-to/go-kit-rest-api-example/services"
)

type Endpoints struct {
	GetTasksEndpoint   endpoint.Endpoint
	GetTaskEndpoint    endpoint.Endpoint
	PostTaskEndpoint   endpoint.Endpoint
	PutTaskEndpoint    endpoint.Endpoint
	DeleteTaskEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s services.Service) Endpoints {
	return Endpoints{
		GetTasksEndpoint:   MakeGetTasksEndpointEndpoint(s),
		GetTaskEndpoint:    MakeGetTaskEndpointEndpoint(s),
		PostTaskEndpoint:   MakePostTaskEndpoint(s),
		PutTaskEndpoint:    MakePutTaskEndpoint(s),
		DeleteTaskEndpoint: MakeDeleteTaskEndpoint(s),
	}
}

func MakeGetTasksEndpointEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		t, e := s.GetTasks(ctx)
		return t, e
	}
}

func MakeGetTaskEndpointEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTaskRequest)
		t, e := s.GetTask(ctx, req.ID)
		return t, e
	}
}

func MakePostTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postTaskRequest)
		t, e := s.PostTask(ctx, req.Task)
		return t, e
	}
}

func MakePutTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putTaskRequest)
		t, e := s.PutTask(ctx, req.ID, req.Task)
		return t, e
	}
}

func MakeDeleteTaskEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteTaskRequest)
		e := s.DeleteTask(ctx, req.ID)
		return nil, e
	}
}

type getTaskRequest struct {
	ID string
}

type postTaskRequest struct {
	Task entities.Task
}

type putTaskRequest struct {
	ID   string
	Task entities.Task
}

type deleteTaskRequest struct {
	ID   string
	Task entities.Task
}
