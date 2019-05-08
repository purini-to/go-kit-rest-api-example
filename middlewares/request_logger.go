package middlewares

import (
	cw "github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

type Logger interface {
	Log(keyvals ...interface{}) error
}

// RequestLogger prodive request log middleware.
func RequestLogger(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := cw.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				keyvals := []interface{}{
					"method", r.Method,
					"url", r.URL.String(),
					"proto", r.Proto,
					"status", ww.Status(),
					"ip", r.RemoteAddr,
					"byte", ww.BytesWritten(),
					"latency", time.Since(t1),
				}
				if reqID := r.Context().Value(cw.RequestIDKey); reqID != nil {
					keyvals = append(keyvals, "reqId", reqID)
				}
				logger.Log(keyvals...)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
