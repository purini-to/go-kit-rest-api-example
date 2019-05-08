package middlewares

import (
	"net/http"
	"runtime/debug"
)

type ErrorFunc func(w http.ResponseWriter, r *http.Request, panicErr interface{})

func Recoverer(fn ErrorFunc, logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logger.Log("panicErr", rvr, "stack", debug.Stack())
					fn(w, r, rvr)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
