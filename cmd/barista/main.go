package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/italolelis/coffee-shop/internal/app/coffees"
	"github.com/italolelis/coffee-shop/internal/app/config"
	"github.com/italolelis/coffee-shop/internal/app/staff"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/italolelis/coffee-shop/internal/pkg/metric"
	"github.com/italolelis/coffee-shop/internal/pkg/stream"
	"github.com/rafaeljesus/rabbus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
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

	startServer(ctx, cfg)

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

func startServer(ctx context.Context, cfg *config.Config) {
	logger := log.WithContext(ctx)

	hc, err := setupHealthCheckers(cfg.EventStream.DSN)
	if err != nil {
		logger.Fatalw("could not start health checks", "err", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to barista service"))
	})
	r.Handle("/metrics", metric.Setup(ctx, cfg.Metrics))
	r.Handle("/status", handlers.NewJSONHandlerFunc(hc, nil))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler: chi.ServerBaseContext(ctx, &ochttp.Handler{
			Handler:     r,
			Propagation: &b3.HTTPFormat{},
		}),
	}

	go func() {
		logger.Infow("service running", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("started web server on ")
		}
	}()
}

func setupHealthCheckers(amqpDSN string) (*health.Health, error) {
	// Create a new health instance
	h := health.New()
	h.DisableLogging()

	if err := h.AddChecks([]*health.Config{
		{
			Name:     "amqp-barista-check",
			Checker:  stream.NewChecker(stream.WithDSN(amqpDSN)),
			Interval: time.Duration(3) * time.Second,
			Fatal:    true,
		},
	}); err != nil {
		return nil, err
	}

	if err := h.Start(); err != nil {
		return nil, err
	}

	return h, nil
}
