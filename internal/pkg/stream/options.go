package stream

import "time"

// Option represents the amqp checker options
type Option func(*AMQPChecker)

// WithDSN sets the amqp dsn
func WithDSN(dsn string) Option {
	return func(c *AMQPChecker) {
		c.DSN = dsn
	}
}

// WithExchange sets the amqp exchange
func WithExchange(ex string) Option {
	return func(c *AMQPChecker) {
		c.Exchange = ex
	}
}

// WithRoutingKey sets the routing key
func WithRoutingKey(rk string) Option {
	return func(c *AMQPChecker) {
		c.RoutingKey = rk
	}
}

// WithQueue sets the queue
func WithQueue(q string) Option {
	return func(c *AMQPChecker) {
		c.Queue = q
	}
}

// WithConsumeTimeout sets the consume timeout
func WithConsumeTimeout(t time.Duration) Option {
	return func(c *AMQPChecker) {
		c.ConsumeTimeout = t
	}
}
