package metric

import (
	"context"
	"net/http"

	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

// Setup sets up the application metrics
func Setup(ctx context.Context, cfg Metrics) http.Handler {
	logger := log.WithContext(ctx)

	if err := view.Register(
		ochttp.ClientCompletedCount,
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
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
