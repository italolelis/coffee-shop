package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/httptrace"
)

// Tracing middleware setting a value on the request context
func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := WithContext(r.Context())
		attrs, entries, spanCtx := httptrace.Extract(r.Context(), r)

		r = r.WithContext(correlation.ContextWithMap(r.Context(), correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		})))

		ctx, span := tr.Start(
			trace.ContextWithRemoteSpanContext(r.Context(), spanCtx),
			r.URL.String(),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
