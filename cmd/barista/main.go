package main

import (
	"context"

	"github.com/italolelis/coffee-shop/internal/config"
	"github.com/italolelis/coffee-shop/internal/log"
	"github.com/italolelis/coffee-shop/internal/stream"
	"github.com/italolelis/coffee-shop/pkg/coffees"
	"github.com/italolelis/coffee-shop/pkg/staff"
	"github.com/rafaeljesus/rabbus"
)

var workforce = []*staff.Barista{
	{
		Name: "Thomas",
		Skills: []coffees.CoffeeType{
			&coffees.Espresso{},
			&coffees.Cappuccino{},
		},
	},
	{
		Name: "Sofia",
		Skills: []coffees.CoffeeType{
			&coffees.Espresso{},
			&coffees.Cappuccino{},
			&coffees.Latte{},
		},
	},
	{
		Name: "John",
		Skills: []coffees.CoffeeType{
			&coffees.Espresso{},
			&coffees.Cappuccino{},
			&coffees.Latte{},
		},
	},
}

func main() {
	// creates a cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// gets the contextual logging
	logger := log.WithContext(ctx)
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal(err)
		}
	}()

	// loads the configuration from the environment
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

	// Setup buffered input/output queues for the workers
	results := make(chan *staff.OrderDone, 512)

	for _, b := range workforce {
		go b.Prepare(ctx, messages, results)
	}

	logger.Debug("listening to orders")
	for {
		o, ok := <-results
		if !ok {
			logger.Debug("stop listening done orders!")
			break
		}

		logger.Infof(
			"%s -> %s size %s for %s your order is ready!",
			o.DoneBy.Name,
			o.Type,
			o.Size,
			o.CustomerName,
		)
	}
}
