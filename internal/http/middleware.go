package http

import (
	"net/http"
	"strconv"
	"time"
)

type timingResponseWriter struct {
	http.ResponseWriter
	start time.Time
}

func (tw *timingResponseWriter) WriteHeader(statusCode int) {
	elapsed := time.Since(tw.start).Microseconds()
	tw.Header().Set(
		"X-Processing-Time-Micros",
		strconv.FormatInt(elapsed, 10),
	)

	tw.ResponseWriter.WriteHeader(statusCode)
}

// ProcessingTimeMiddleware adds X-Processing-Time-Micros header
// to all HTTP responses
func ProcessingTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := &timingResponseWriter{
			ResponseWriter: w,
			start:          time.Now(),
		}

		next.ServeHTTP(tw, r)
	})
}
