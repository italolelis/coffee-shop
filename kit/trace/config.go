package trace

type Tracing struct {
	ServiceName       string  `envconfig:"JAEGER_SERVICE_NAME"`
	CollectorEndpoint string  `envconfig:"JAEGER_COLLECTOR_ENDPOINT"`
	ProbabilityFactor float64 `envconfig:"JAEGER_PROBABILITY_FACTOR" default="1"`
}
