package go_kit_rest_api_example

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/purini-to/go-kit-rest-api-example/services"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	GetTasksEndpoint endpoint.Endpoint
	GetTaskEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s services.Service) Endpoints {
	return Endpoints{
		GetTasksEndpoint: MakeGetTasksEndpointEndpoint(s),
		GetTaskEndpoint:  MakeGetTaskEndpointEndpoint(s),
	}
}

// MakeGetTasksEndpointEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetTasksEndpointEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		t, e := s.GetTasks(ctx)
		return t, e
	}
}

// MakeGetTaskEndpointEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetTaskEndpointEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTaskRequest)
		t, e := s.GetTask(ctx, req.ID)
		return t, e
	}
}

type getTaskRequest struct {
	ID string
}
