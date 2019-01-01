package stream

import "time"

// EventStream holds the event stream configurations
type EventStream struct {
	DSN       string        `envconfig:"EVENT_STREAM_DSN"`
	Attempts  int           `envconfig:"EVENT_STREAM_RETRY_ATTEMPTS" default:"5"`
	Backoff   time.Duration `envconfig:"EVENT_STREAM_RETRY_BACKOFF" default:"2s"`
	Threshold uint32        `envconfig:"EVENT_STREAM_RETRY_THRESHOLD" default:"5"`
}
