package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/plugin/grpctrace"

	"github.com/italolelis/coffee-shop/internal/app/http/rest"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/italolelis/coffee-shop/internal/pkg/pb"
	"github.com/italolelis/coffee-shop/internal/pkg/signal"
	"github.com/italolelis/coffee-shop/internal/pkg/tracing"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type config struct {
	LogLevel string `split_words:"true" default:"info"`
	API      struct {
		Addr            string        `split_words:"true" default:"0.0.0.0:8080"`
		ReadTimeout     time.Duration `split_words:"true" default:"30s"`
		WriteTimeout    time.Duration `split_words:"true" default:"30s"`
		IdleTimeout     time.Duration `split_words:"true" default:"5s"`
		ShutdownTimeout time.Duration `split_words:"true" default:"5s"`
	}
	Payment struct {
		Addr    string        `split_words:"true" required:"true"`
		Timeout time.Duration `split_words:"true" default:"2s"`
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
	// Start REST Service
	// =========================================================================
	logger.Infow("diling payment service", "addr", cfg.Payment.Addr)
	paymentDiler, err := grpc.DialContext(
		ctx,
		cfg.Payment.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Payment.Timeout),
		grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor(t)),
		grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(t)),
	)
	if err != nil {
		return fmt.Errorf("failed to dial payment service: %w", err)
	}
	defer paymentDiler.Close()

	s := rest.NewServer(
		rest.Config{
			Addr:         cfg.API.Addr,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
			IdleTimeout:  cfg.API.IdleTimeout,
		},
		pb.NewPaymentClient(paymentDiler),
	)
	go func() {
		logger.Infow("Initializing REST support", "host", cfg.API.Addr)
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

		ctx, cancel := context.WithTimeout(ctx, cfg.API.ShutdownTimeout)
		defer cancel()

		if err := s.Stop(ctx); err != nil {
			return fmt.Errorf("failed to close probe server gracefully: %w", err)
		}
	}

	return nil
}
