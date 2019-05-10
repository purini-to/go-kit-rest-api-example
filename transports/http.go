package transports

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/purini-to/go-kit-rest-api-example/endpoints"
	"github.com/purini-to/go-kit-rest-api-example/entities"
	"github.com/purini-to/go-kit-rest-api-example/services"
	"go.uber.org/zap"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
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

type requestLogger struct {
	wrapZapLogger
}

func (e *requestLogger) Log(keyvals ...interface{}) error {
	logger := e.With(keyvals...)
	logger.Info("request")
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

func requestLoggerMiddle(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func(begin time.Time) {
				keyvals := []interface{}{
					"method", r.Method,
					"url", r.URL.String(),
					"proto", r.Proto,
					"status", ww.Status(),
					"ip", r.RemoteAddr,
					"byte", ww.BytesWritten(),
					"latency", time.Since(begin),
				}
				if reqID := middleware.GetReqID(r.Context()); len(reqID) > 0 {
					keyvals = append([]interface{}{"reqId", reqID}, keyvals...)
				}
				logger.Log(keyvals...)
			}(time.Now())
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

type errorFunc func(w http.ResponseWriter, r *http.Request, panicErr interface{})

func recovererMiddle(fn errorFunc, logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					keyvals := []interface{}{"panicErr", rvr}
					if reqID := middleware.GetReqID(r.Context()); len(reqID) > 0 {
						keyvals = append([]interface{}{"reqId", reqID}, keyvals...)
					}
					logger.Log(keyvals...)
					fn(w, r, rvr)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func MakeHTTPHandler(s services.Service, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	// A good base middleware stack
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		requestLoggerMiddle(&requestLogger{wrapZapLogger: wrapZapLogger{logger: logger}}),
		recovererMiddle(panicHandler, &recoverLogger{wrapZapLogger: wrapZapLogger{logger: logger}}),
	)

	r.Route("/tasks", func(r chi.Router) {

		r.Method("GET", "/", httptransport.NewServer(
			endpoints.GetTasksEndpoint(s),
			decodeEmptyRequest,
			encodeGetTasksResponse,
			options...,
		))

		r.Method("POST", "/", httptransport.NewServer(
			endpoints.PostTaskEndpoint(s),
			decodePostTaskRequest,
			encodePostResponse,
			options...,
		))

		r.Route("/{id}", func(r chi.Router) {

			r.Method("GET", "/", httptransport.NewServer(
				endpoints.GetTaskEndpoint(s),
				decodeGetTaskRequest,
				encodeGetTaskResponse,
				options...,
			))

			r.Method("PUT", "/", httptransport.NewServer(
				endpoints.PutTaskEndpoint(s),
				decodePutTaskRequest,
				encodePutResponse,
				options...,
			))

			r.Method("DELETE", "/", httptransport.NewServer(
				endpoints.DeleteTaskEndpoint(s),
				decodeDeleteProfileRequest,
				encodeDeleteResponse,
				options...,
			))

		})
	})

	return r
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodePostTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.PostTaskInput
	if e := json.NewDecoder(r.Body).Decode(&req.Task); e != nil {
		return nil, e
	}

	return req, nil
}

func decodeGetTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")

	return endpoints.GetTaskInput{ID: id}, nil
}

func decodePutTaskRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")

	var task entities.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, err
	}
	return endpoints.PutTaskInput{
		ID:   id,
		Task: task,
	}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")

	return endpoints.DeleteTaskInput{ID: id}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeEmptyResponse(_ context.Context, w http.ResponseWriter, _ interface{}) error {
	_, err := w.Write([]byte(""))
	return err
}

func encodeGetTasksResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(endpoints.GetTasksOutput)
	return encodeResponse(ctx, w, r.Tasks)
}

func encodeGetTaskResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(endpoints.GetTaskOutput)
	return encodeResponse(ctx, w, r.Task)
}

func encodePostResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusCreated)
	return encodeEmptyResponse(ctx, w, response)
}

func encodePutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return encodeEmptyResponse(ctx, w, response)
}

func encodeDeleteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return encodeEmptyResponse(ctx, w, response)
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

func panicHandler(w http.ResponseWriter, r *http.Request, _ interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": http.StatusText(http.StatusInternalServerError),
	})
}

func codeFrom(err error) int {
	switch err {
	case services.ErrNotFound:
		return http.StatusNotFound
	case services.ErrInconsistentIDs:
		return http.StatusBadRequest
	case services.ErrAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
