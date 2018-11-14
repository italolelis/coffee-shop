package main

import (
	"context"

	"github.com/italolelis/barista/pkg/config"
	"github.com/golang/protobuf/proto"
	"github.com/italolelis/kit/proto/order"
	"github.com/italolelis/kit/log"
	"github.com/rafaeljesus/rabbus"
)

func main() {
	// creates a cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// gets the contextual logging
	logger := log.WithContext(ctx)
	defer logger.Sync()

	// loads the configuration from the enviroment
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	log.SetLevel(cfg.LogLevel)

	// setup the event stream. In this case is an event broker because we chose rabbitmq
	eventStream, err := setupEventStream(ctx, cfg.EventStream)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func(r *rabbus.Rabbus) {
		if err := r.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(eventStream)

	go eventStream.Run(ctx)

	messages, err := eventStream.Listen(rabbus.ListenConfig{
		Exchange: "orders",
		Kind:     "topic",
		Key:      "orders.created",
		Queue:    "orders_barista",
	})
	if err != nil {
		logger.Fatalw("failed to create listener", "err", err.Error())
		return
	}
	defer close(messages)

	logger.Info("listening for messages...")
	for {
		m, ok := <-messages
		if !ok {
			logger.Info("stop listening messages!")
			break
		}

		o := order.Created{}
		err = proto.Unmarshal(m.Body, &o)
		if err != nil {
			logger.Errorw("unmarshaling error", "err", err)
		}

		logger.Infow("your order is ready", "order_id", o.ID)
		m.Ack(false)

		logger.Info("message was consumed")
	}
}

func setupEventStream(ctx context.Context, cfg config.EventStream) (*rabbus.Rabbus, error) {
	logger := log.WithContext(ctx)

	cbStateChangeFunc := func(name, from, to string) {
		logger.Debugw("rabbitmq state changed", "from", from, "to", to)
	}

	return rabbus.New(
		cfg.DSN,
		rabbus.Durable(true),
		rabbus.Attempts(cfg.Attempts),
		rabbus.Sleep(cfg.Backoff),
		rabbus.Threshold(cfg.Threshold),
		rabbus.OnStateChange(cbStateChangeFunc),
	)
}
