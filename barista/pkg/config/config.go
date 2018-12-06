package config

import (
	"github.com/italolelis/kit/stream"
	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration
type Config struct {
	LogLevel    string `envconfig:"LOG_LEVEL"`
	EventStream stream.EventStream
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
