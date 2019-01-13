package stream

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"time"
)

const (
	defaultExchange       = "health_check"
	defaultConsumeTimeout = time.Second * 3
)

// AMQPChecker is the health checker for amqp
type AMQPChecker struct {
	// DSN is the RabbitMQ instance connection DSN. Required.
	DSN string
	// Exchange is the application health check exchange. If not set - "health_check" is used.
	Exchange string
	// RoutingKey is the application health check routing key within health check exchange.
	// Can be an application or host name, for example.
	// If not set - host name is used.
	RoutingKey string
	// Queue is the application health check queue, that binds to the exchange with the routing key.
	// If not set - "<exchange>.<routing-key>" is used.
	Queue string
	// ConsumeTimeout is the duration that health check will try to consume published test message.
	// If not set - 3 seconds
	ConsumeTimeout time.Duration
}

// NewChecker creates a new instance of AMQPChecker
func NewChecker(opts ...Option) *AMQPChecker {
	h := AMQPChecker{
		Exchange:       defaultExchange,
		ConsumeTimeout: defaultConsumeTimeout,
	}

	for _, opt := range opts {
		opt(&h)
	}

	if h.RoutingKey == "" {
		host, err := os.Hostname()
		if nil != err {
			h.RoutingKey = "-unknown-"
		}
		h.RoutingKey = host
	}

	if h.Queue == "" {
		h.Queue = fmt.Sprintf("%s.%s", h.Exchange, h.RoutingKey)
	}

	return &h
}

// Status is used for performing an HTTP check against a dependency; it satisfies
// the "ICheckable" interface.
func (h *AMQPChecker) Status() (interface{}, error) {
	conn, err := amqp.Dial(h.DSN)
	if err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed on dial phase: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed on getting channel phase: %s", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(h.Exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed during declaring exchange: %s", err)
	}

	if _, err := ch.QueueDeclare(h.Queue, false, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed during declaring queue: %s", err)
	}

	if err := ch.QueueBind(h.Queue, h.RoutingKey, h.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed during binding: %s", err)
	}

	messages, err := ch.Consume(h.Queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed during consuming: %s", err)
	}

	done := make(chan struct{})

	go func() {
		// block until: a message is received, or message channel is closed (consume timeout)
		<-messages

		// release the channel resources, and unblock the receive on done below
		close(done)

		// now drain any incidental remaining messages
		for range messages {
		}
	}()

	p := amqp.Publishing{Body: []byte(time.Now().Format(time.RFC3339Nano))}
	if err := ch.Publish(h.Exchange, h.RoutingKey, false, false, p); err != nil {
		return nil, fmt.Errorf("rabbitMQ health check failed during publishing: %s", err)
	}

	return nil, nil
}
