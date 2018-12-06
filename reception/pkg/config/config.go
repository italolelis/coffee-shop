package config

import (
	"github.com/italolelis/kit/metric"
	"github.com/italolelis/kit/stream"
	"github.com/italolelis/kit/trace"
	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration
type Config struct {
	Port        int    `default:"8080"`
	LogLevel    string `envconfig:"LOG_LEVEL"`
	Database    Database
	EventStream stream.EventStream
	Tracing     trace.Tracing
	Metrics     metric.Metrics
}

// Database holds the database configurations
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
}

// Load loads the app config from the environment
func Load() (*Config, error) {
	var s Config
	err := envconfig.Process("", &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
