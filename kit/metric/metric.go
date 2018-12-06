package metric

import (
	"context"
	"net/http"

	"github.com/italolelis/kit/log"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

// Setup sets up the application metrics
func Setup(ctx context.Context, cfg Metrics) http.Handler {
	logger := log.WithContext(ctx)

	if err := view.Register(
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
	); err != nil {
		logger.Fatal(err)
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: cfg.ServiceName,
	})
	if err != nil {
		logger.Fatalw("failed to create the prometheus stats exporter", "err", err.Error())
	}
	view.RegisterExporter(exporter)
	view.SetReportingPeriod(cfg.ReportingPeriod)

	return exporter
}
