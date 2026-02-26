package middleware

import (
	"net/http"
	"time"

	"knowledge-capsule/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.InfoRequest(r, rw.status, time.Since(start))
	})
}
