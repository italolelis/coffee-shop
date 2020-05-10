package main

import (
	"context"
	"fmt"
	"time"

	"github.com/italolelis/coffee-shop/internal/app/http/grpc"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/italolelis/coffee-shop/internal/pkg/signal"
	"github.com/italolelis/coffee-shop/internal/pkg/tracing"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	LogLevel string `split_words:"true" default:"info"`
	Web      struct {
		Addr            string        `split_words:"true" default:"0.0.0.0:8081"`
		ShutdownTimeout time.Duration `split_words:"true" default:"5s"`
	}
	Tracing struct {
		Addr        string `split_words:"true"`
		ServiceName string `split_words:"true"`
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.WithContext(ctx)
	logger.Sync()

	if err := run(ctx); err != nil {
		logger.Fatalw("an error happened", "err", err)
	}
}

func run(ctx context.Context) error {
	logger := log.WithContext(ctx)

	// =========================================================================
	// Configuration
	// =========================================================================
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load the env vars: %w", err)
	}

	log.SetLevel(cfg.LogLevel)

	// =========================================================================
	// Start Tracing Support
	// =========================================================================
	logger.Info("initializing tracing support")
	tp, close, err := tracing.InitTracer(cfg.Tracing.Addr, cfg.Tracing.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to setup tracing: %w", err)
	}
	defer close()

	t := tp.Tracer("main")
	ctx = tracing.NewContext(ctx, t)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	var serverErrors = make(chan error, 1)

	// =========================================================================
	// Start GRPC Service
	// =========================================================================
	s := grpc.NewServer(grpc.Config{Addr: cfg.Web.Addr}, tp.Tracer("main"))
	go func() {
		logger.Infow("Initializing GRPC support", "addr", cfg.Web.Addr)
		serverErrors <- s.ListenAndServe(ctx)
	}()

	// =========================================================================
	// Signal notifier
	// =========================================================================
	done := signal.New(ctx)

	logger.Info("application running")

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		logger.Info("shutdown")

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		s.Stop(ctx)
	}

	return nil
}
