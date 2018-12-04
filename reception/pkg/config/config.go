package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration
type Config struct {
	Port        int    `default:"8080"`
	LogLevel    string `envconfig:"LOG_LEVEL"`
	Database    Database
	EventStream EventStream
	Tracing     Tracing
}

// Database holds the database configurations
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
}

// EventStream holds the event stream configurations
type EventStream struct {
	DSN       string        `envconfig:"EVENT_STREAM_DSN"`
	Attempts  int           `envconfig:"EVENT_STREAM_RETRY_ATTEMPTS" default:"5"`
	Backoff   time.Duration `envconfig:"EVENT_STREAM_RETRY_BACKOFF" default:"2s"`
	Threshold uint32        `envconfig:"EVENT_STREAM_RETRY_THRESHOLD" default:"5"`
}

type Tracing struct {
	ServiceName       string  `envconfig:"JAEGER_SERVICE_NAME"`
	CollectorEndpoint string  `envconfig:"JAEGER_COLLECTOR_ENDPOINT"`
	ProbabilityFactor float64 `envconfig:"JAEGER_PROBABILITY_FACTOR" default="1"`
}

// Load loads the app config from the enviroment
func Load() (*Config, error) {
	var s Config
	err := envconfig.Process("", &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
