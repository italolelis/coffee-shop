package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration
type Config struct {
	LogLevel    string `envconfig:"LOG_LEVEL"`
	EventStream EventStream
}

// EventStream holds the event stream configurations
type EventStream struct {
	DSN       string        `envconfig:"EVENT_STREAM_DSN"`
	Attempts  int           `envconfig:"EVENT_STREAM_RETRY_ATTEMPTS" default:"5"`
	Backoff   time.Duration `envconfig:"EVENT_STREAM_RETRY_BACKOFF" default:"2s"`
	Threshold uint32        `envconfig:"EVENT_STREAM_RETRY_THRESHOLD" default:"5"`
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
