package tracing

import (
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer creates a new trace provider instance.
func InitTracer(addr string, serviceName string) (*sdktrace.Provider, func(), error) {
	// Create and install Jaeger export pipeline
	return jaeger.NewExportPipeline(
		jaeger.WithCollectorEndpoint(addr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
}

