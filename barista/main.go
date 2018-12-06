package main

import (
	"context"
	"github.com/italolelis/kit/stream"

	"github.com/golang/protobuf/proto"
	"github.com/italolelis/barista/pkg/config"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/kit/proto/order"
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

	eventStream, flush := stream.Setup(ctx, cfg.EventStream)
	defer flush()

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

	logger.Debug("listening for messages...")
	for {
		m, ok := <-messages
		if !ok {
			logger.Debug("stop listening messages!")
			break
		}

		o := order.Created{}
		err = proto.Unmarshal(m.Body, &o)
		if err != nil {
			logger.Errorw("unmarshaling error", "err", err)
		}

		for _, i := range o.Items {
			logger.Infof("%s size %s for %s your order is ready!", i.Type, i.Size, o.CustomerName)
		}

		m.Ack(false)

		logger.Debug("message was consumed")
	}
}
