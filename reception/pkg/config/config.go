package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration
type Config struct {
	Port     int    `default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	Database Database
}

// Database holds the database configurations
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
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
