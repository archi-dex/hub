package api

import (
	"net/http"

	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/mux"
)

func loggingMiddlware(logger util.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logger.Infow("request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
			)

			next.ServeHTTP(rw, r)
		})
	}
}

func defaultHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(rw, r)
	})
}
