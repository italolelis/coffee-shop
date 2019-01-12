package stream

import (
	"context"

	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/rafaeljesus/rabbus"
)

// Setup sets up the event stream. In this case is an event broker because we chose rabbitmq
func Setup(ctx context.Context, cfg EventStream) (*rabbus.Rabbus, func()) {
	logger := log.WithContext(ctx)

	cbStateChangeFunc := func(name, from, to string) {
		logger.Debugw("rabbitmq state changed", "from", from, "to", to)
	}

	eventStream, err := rabbus.New(
		cfg.DSN,
		rabbus.Durable(true),
		rabbus.Attempts(cfg.Attempts),
		rabbus.Sleep(cfg.Backoff),
		rabbus.Threshold(cfg.Threshold),
		rabbus.OnStateChange(cbStateChangeFunc),
	)
	if err != nil {
		logger.Fatalw("failed to establish the rabbitmq connection", "err", err.Error())
	}

	go func() {
		for {
			select {
			case <-eventStream.EmitOk():
				logger.Debug("message sent")
			case <-eventStream.EmitErr():
				logger.Debug("message was not sent")
			}
		}
	}()

	go func() {
		if err := eventStream.Run(ctx); err != nil {
			logger.Fatalw("failed to initialize rabbitmq channels", "err", err.Error())
		}
	}()

	return eventStream, func() {
		if err := eventStream.Close(); err != nil {
			logger.Errorw("failed to close rabbitmq connection", "err", err.Error())
		}
	}
}
