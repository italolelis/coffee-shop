package log

import (
	"net/http"
	"net/url"
	"time"

	"github.com/felixge/httpsnoop"
)

// NewMiddleware creates a new log middleware
func NewMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := WithContext(r.Context())

		logger.Infow("Started request",
			"method", r.Method,
			"path", r.URL.Path,
		)

		// reverse proxy replaces original request with target request, so keep original one
		originalURL := &url.URL{}
		*originalURL = *r.URL

		logger = logger.With(
			"method", r.Method,
			"host", r.Host,
			"request", r.RequestURI,
			"remote-addr", r.RemoteAddr,
			"referer", r.Referer(),
			"user-agent", r.UserAgent(),
		)

		m := httpsnoop.CaptureMetrics(handler, w, r)

		logger = logger.With(
			"code", m.Code,
			"duration", int(m.Duration/time.Millisecond),
			"duration-fmt", m.Duration.String(),
		)

		logger.Info("Completed handling request")
	})
}
