package go_kit_rest_api_example

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/purini-to/go-kit-rest-api-example/middlewares"
	"go.uber.org/zap"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

type wrapZapLogger struct {
	logger *zap.Logger
}

func (w *wrapZapLogger) With(keyvals ...interface{}) *zap.Logger {
	logger := w.logger
	for i := 0; i < len(keyvals); i = i + 2 {
		logger = logger.With(zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	return logger
}

type errorLogger struct {
	wrapZapLogger
}

func (e *errorLogger) Log(keyvals ...interface{}) error {
	logger := e.With(keyvals...)
	logger.Error("http transport server error")
	return nil
}

type requestLogger struct {
	wrapZapLogger
}

func (e *requestLogger) Log(keyvals ...interface{}) error {
	logger := e.With(keyvals...)
	logger.Info("request log")
	return nil
}

type recoverLogger struct {
	wrapZapLogger
}

func (e *recoverLogger) Log(keyvals ...interface{}) error {
	logger := e.With(keyvals...)
	logger.Error("http transport panic error")
	return nil
}

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(s Service, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(&errorLogger{wrapZapLogger: wrapZapLogger{logger: logger}}),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewares.RequestLogger(&requestLogger{wrapZapLogger: wrapZapLogger{logger: logger}}))
	r.Use(middlewares.Recoverer(func(w http.ResponseWriter, r *http.Request, panicErr interface{}) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": http.StatusText(http.StatusInternalServerError),
		})
	}, &recoverLogger{wrapZapLogger: wrapZapLogger{logger: logger}}))

	r.Method("GET", "/tasks", httptransport.NewServer(
		e.GetTasksEndpoint,
		decodeEmptyRequest,
		encodeResponse,
		options...,
	))

	r.Method("GET", "/tasks/{id}", httptransport.NewServer(
		e.GetTaskEndpoint,
		decodeGetTaskRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeGetTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		return nil, ErrBadRouting
	}

	return getTaskRequest{ID: id}, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
