package trace

import (
	"context"

	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// Setup register the Jaeger exporter to be able to retrieve
// the collected spans.
func Setup(ctx context.Context, cfg Tracing) func() {
	logger := log.WithContext(ctx)

	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: cfg.CollectorEndpoint,
		Process: jaeger.Process{
			ServiceName: cfg.ServiceName,
		},
	})
	if err != nil {
		logger.Errorw("could not create the jaeger exporter", "err", err)
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() { exporter.Flush() }
}
