package config

import (
	"github.com/italolelis/coffee-shop/internal/pkg/metric"
	"github.com/italolelis/coffee-shop/internal/pkg/stream"
	"github.com/italolelis/coffee-shop/internal/pkg/trace"
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config is the application configuration
type Config struct {
	Port              int           `default:"8080"`
	LogLevel          string        `envconfig:"LOG_LEVEL"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT"  default:"10s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT"  default:"120s"`
	Database          Database
	EventStream       stream.EventStream
	Tracing           trace.Tracing
	Metrics           metric.Metrics
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
