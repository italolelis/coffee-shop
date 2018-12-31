package metric

import "time"

type Metrics struct {
	ServiceName     string        `envconfig:"METRICS_SERVICE_NAME"`
	ReportingPeriod time.Duration `envconfig:"METRICS_REPORTING_PERIOD" default:"1s"`
}
