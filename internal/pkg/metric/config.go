package metric

import "time"

// Metrics are the configurations for metrics
type Metrics struct {
	ServiceName     string        `envconfig:"METRICS_SERVICE_NAME"`
	ReportingPeriod time.Duration `envconfig:"METRICS_REPORTING_PERIOD" default:"1s"`
}
