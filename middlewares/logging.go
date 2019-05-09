package middlewares

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/purini-to/go-kit-rest-api-example/entities"
	"github.com/purini-to/go-kit-rest-api-example/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ctxKeyLogger int

const LoggerKey ctxKeyLogger = 0

type Service func(task services.Service) services.Service

func Logging(logger *zap.Logger) Service {
	return func(next services.Service) services.Service {
		return &logging{
			next:   next,
			logger: logger,
		}
	}
}

type logging struct {
	next   services.Service
	logger *zap.Logger
}

func (l logging) PostTask(ctx context.Context, t entities.Task) (task entities.Task, err error) {
	ctx = withLogger(ctx, l.logger)
	defer func(begin time.Time) {
		write(ctx, "PostTask", err, begin, zap.String("id", t.ID))
	}(time.Now())
	return l.next.PostTask(ctx, t)
}

func (l logging) PutTask(ctx context.Context, id string, t entities.Task) (task entities.Task, err error) {
	ctx = withLogger(ctx, l.logger)
	defer func(begin time.Time) {
		write(ctx, "PutTask", err, begin, zap.String("id", id))
	}(time.Now())
	return l.next.PutTask(ctx, id, t)
}

func (l logging) DeleteTask(ctx context.Context, id string) (err error) {
	ctx = withLogger(ctx, l.logger)
	defer func(begin time.Time) {
		write(ctx, "DeleteTask", err, begin, zap.String("id", id))
	}(time.Now())
	return l.next.DeleteTask(ctx, id)
}

func (l logging) GetTasks(ctx context.Context) (tasks []entities.Task, err error) {
	ctx = withLogger(ctx, l.logger)
	defer func(begin time.Time) {
		write(ctx, "GetTasks", err, begin)
	}(time.Now())
	return l.next.GetTasks(ctx)
}

func (l logging) GetTask(ctx context.Context, id string) (task entities.Task, err error) {
	ctx = withLogger(ctx, l.logger)
	defer func(begin time.Time) {
		write(ctx, "GetTask", err, begin, zap.String("id", id))
	}(time.Now())
	return l.next.GetTask(ctx, id)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return nil
	}
	if logger, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return logger
	}
	return nil
}

func withLogger(ctx context.Context, logger *zap.Logger) context.Context {
	l := logger
	if reqID := middleware.GetReqID(ctx); len(reqID) > 0 {
		l = logger.With(zap.String("reqId", reqID))
	}
	return context.WithValue(ctx, LoggerKey, l)
}

func write(ctx context.Context, method string, err error, begin time.Time, params ...zap.Field) {
	logger := GetLogger(ctx)
	if logger == nil {
		return
	}

	fields := []zap.Field{
		zap.String("method", method),
		zap.Error(err),
		zap.Duration("latency", time.Since(begin)),
	}
	if len(params) > 0 {
		fields = append(fields, params...)
	}

	ch := logger.Check(levelFrom(err), "service")
	ch.Write(fields...)
}

func levelFrom(err error) zapcore.Level {
	switch err {
	case nil, services.ErrNotFound, services.ErrAlreadyExists, services.ErrInconsistentIDs:
		return zap.InfoLevel
	default:
		return zap.ErrorLevel
	}
}
