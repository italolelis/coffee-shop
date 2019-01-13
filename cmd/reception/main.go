package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/checkers"
	"github.com/InVisionApp/go-health/handlers"
	"github.com/go-chi/chi"
	"github.com/italolelis/coffee-shop/internal/app/coffees"
	"github.com/italolelis/coffee-shop/internal/app/config"
	"github.com/italolelis/coffee-shop/internal/app/http/rest"
	"github.com/italolelis/coffee-shop/internal/app/order"
	"github.com/italolelis/coffee-shop/internal/app/storage/postgres"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/italolelis/coffee-shop/internal/pkg/metric"
	"github.com/italolelis/coffee-shop/internal/pkg/stream"
	"github.com/italolelis/coffee-shop/internal/pkg/trace"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

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

	if err := run(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}

func run(ctx context.Context) error {
	// gets the contextual logging
	logger := log.WithContext(ctx)

	// loads the configuration from the environment
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log.SetLevel(cfg.LogLevel)

	db, dbClose := setupDatabase(ctx, cfg.Database)
	defer dbClose()

	eventStream, streamClose := stream.Setup(ctx, cfg.EventStream)
	defer streamClose()

	flush := trace.Setup(ctx, cfg.Tracing)
	defer flush()

	metricsHandler := metric.Setup(ctx, cfg.Metrics)

	hc, err := setupHealthCheckers(db.DB, cfg.EventStream.DSN)
	if err != nil {
		return err
	}

	// coffee setup
	cwr := postgres.NewCoffeeWriteRepository(db)
	crr := postgres.NewCoffeeReadRepository(db)
	cs := coffees.NewService(cwr, crr)

	// order setup
	owr := postgres.NewPostgresOrderWriteRepository(db)
	orr := postgres.NewPostgresOrderReadRepository(db)
	os := order.NewService(owr, orr, crr, eventStream)

	logger.Infow("service running", "port", cfg.Port)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler: chi.ServerBaseContext(ctx, &ochttp.Handler{
			Handler:     rest.NewServer(cs, os, metricsHandler, handlers.NewJSONHandlerFunc(hc, nil)),
			Propagation: &b3.HTTPFormat{},
		}),
	}

	return srv.ListenAndServe()
}

// setupDatabase connects to the primary data store
func setupDatabase(ctx context.Context, cfg config.Database) (*sqlx.DB, func()) {
	logger := log.WithContext(ctx)

	// Register our ocsql wrapper for the provided Postgres driver.
	driverName, err := ocsql.Register("postgres", ocsql.WithAllTraceOptions())
	if err != nil {
		logger.Fatalw("could not register the wrapper provider", "err", err)
	}

	// Connect to a Postgres database using the ocsql driver wrapper.
	rawDB, err := sql.Open(driverName, cfg.DSN)
	if err != nil {
		logger.Fatalw("could not open a connection with the database", "driver", driverName, "err", err)
	}

	// Wrap our *sql.DB with sqlx.
	db := sqlx.NewDb(rawDB, "postgres")
	return db, func() {
		if err := db.Close(); err != nil {
			logger.Fatal(err)
		}
	}
}

func setupHealthCheckers(db *sql.DB, amqpDSN string) (*health.Health, error) {
	// Create a new health instance
	h := health.New()
	h.DisableLogging()

	sqlCheck, err := checkers.NewSQL(&checkers.SQLConfig{
		Pinger: db,
	})
	if err != nil {
		return nil, err
	}

	if err = h.AddChecks([]*health.Config{
		{
			Name:     "db-reception-check",
			Checker:  sqlCheck,
			Interval: time.Duration(3) * time.Second,
			Fatal:    true,
		},
		{
			Name:     "amqp-reception-check",
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
