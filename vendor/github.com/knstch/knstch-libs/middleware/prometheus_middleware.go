package middleware

import (
	"net/http"
	"strconv"
	"time"

	metrics "github.com/knstch/knstch-libs/prometeus"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func WithTrackingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		status := rec.status
		if status == 0 {
			status = 200
		}

		metrics.RequestCount.With("path", path, "code", strconv.Itoa(status)).Add(1)
		metrics.RequestDuration.With("path", path, "code", strconv.Itoa(status)).Observe(duration)
	})
}
