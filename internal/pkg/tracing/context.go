package tracing

import (
	"context"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

type tracingKeyType int

const tracingKey tracingKeyType = iota

func NewContext(ctx context.Context, t trace.Tracer) context.Context {
	return context.WithValue(ctx, tracingKey, t)
}

func WithContext(ctx context.Context) trace.Tracer {
	if ctxTracing, ok := ctx.Value(tracingKey).(trace.Tracer); ok {
		return ctxTracing
	}

	return global.Tracer("main")
}

func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	tr := WithContext(ctx)
	return tr.Start(ctx, name)
}
