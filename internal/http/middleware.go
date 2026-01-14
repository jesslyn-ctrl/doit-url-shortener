package http

import (
	stdhttp "net/http"
	"strconv"
	"time"
)

type timingResponseWriter struct {
	stdhttp.ResponseWriter
	start time.Time
}

func (tw *timingResponseWriter) WriteHeader(statusCode int) {
	elapsed := time.Since(tw.start).Microseconds()
	tw.ResponseWriter.Header().Set(
		"X-Processing-Time-Micros",
		strconv.FormatInt(elapsed, 10),
	)

	tw.ResponseWriter.WriteHeader(statusCode)
}

func (tw *timingResponseWriter) Write(b []byte) (int, error) {
	tw.ResponseWriter.Header().Set(
		"X-Processing-Time-Micros",
		strconv.FormatInt(time.Since(tw.start).Microseconds(), 10),
	)
	return tw.ResponseWriter.Write(b)
}

// ProcessingTimeMiddleware adds X-Processing-Time-Micros header
// to all HTTP responses
func ProcessingTimeMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		tw := &timingResponseWriter{
			ResponseWriter: w,
			start:          time.Now(),
		}

		next.ServeHTTP(tw, r)
	})
}
