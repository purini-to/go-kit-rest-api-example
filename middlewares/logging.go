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

func (l logging) GetTasks(ctx context.Context) (tasks []entities.Task, err error) {
	defer func(begin time.Time) {
		l.log(ctx, "GetTasks", err, begin)
	}(time.Now())
	return l.next.GetTasks(ctx)
}

func (l logging) GetTask(ctx context.Context, id string) (task entities.Task, err error) {
	defer func(begin time.Time) {
		l.log(ctx, "GetTask", err, begin, zap.String("id", id))
	}(time.Now())
	return l.next.GetTask(ctx, id)
}

func (l logging) log(ctx context.Context, method string, err error, begin time.Time, params ...zap.Field) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.Error(err),
		zap.Duration("latency", time.Since(begin)),
	}
	if len(params) > 0 {
		fields = append(fields, params...)
	}
	if reqID := middleware.GetReqID(ctx); len(reqID) > 0 {
		fields = append(fields, zap.String("reqId", reqID))
	}

	ch := l.logger.Check(levelFrom(err), "service")
	ch.Write(fields...)
}

func levelFrom(err error) zapcore.Level {
	switch err {
	case nil:
		return zap.InfoLevel
	case services.ErrNotFound:
		return zap.InfoLevel
	default:
		return zap.ErrorLevel
	}
}