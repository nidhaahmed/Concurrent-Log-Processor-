package api

import (
	"logprocessor/pkg/logger"
	"logprocessor/pkg/metrics"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Info(
			r.Method + " " + r.URL.Path + " took " + duration.String(),
		)

		metrics.RecordRequest(duration)
	})
}