package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/reception/pkg/coffees"
	"github.com/italolelis/reception/pkg/config"
	"github.com/italolelis/reception/pkg/order"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rafaeljesus/rabbus"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
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

	// loads the configuration from the environment
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	log.SetLevel(cfg.LogLevel)

	db, close := setupDatabase(ctx, cfg.Database)
	defer close()

	metricsHandler := setupMetrics(ctx)

	eventStream, close := setupEventStream(ctx, cfg.EventStream)
	defer close()

	flush := setupTracing(ctx, cfg.Tracing)
	defer flush()

	wRepo := order.NewPostgresWriteRepository(db)
	rRepo := order.NewPostgresReadRepository(db)
	coffeeReadRepo := coffees.NewPostgresReadRepository(db)
	orderHandler := order.NewHandler(wRepo, rRepo, coffeeReadRepo, eventStream)

	// creates the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/metrics", metricsHandler)
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/{id}", orderHandler.GetOrder)
	})

	logger.Infow("service running", "port", cfg.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), chi.ServerBaseContext(ctx, &ochttp.Handler{
		Handler:     r,
		Propagation: &b3.HTTPFormat{},
	})))
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

// setupEventStream sets up the event stream. In this case is an event broker because we chose rabbitmq
func setupEventStream(ctx context.Context, cfg config.EventStream) (*rabbus.Rabbus, func()) {
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
		logger.Fatal(err.Error())
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
			logger.Fatal(err)
		}
	}()

	return eventStream, func() {
		if err := eventStream.Close(); err != nil {
			logger.Error(err.Error())
		}
	}
}

// setupTracing Register the Jaeger exporter to be able to retrieve
// the collected spans.
func setupTracing(ctx context.Context, cfg config.Tracing) func() {
	logger := log.WithContext(ctx)

	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: cfg.CollectorEndpoint,
		Process: jaeger.Process{
			ServiceName: cfg.ServiceName,
		},
	})
	if err != nil {
		logger.Errorw("could not create the jaeger exporter", "err", err)
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() { exporter.Flush() }
}

// setupMetrics sets up the application metrics
func setupMetrics(ctx context.Context) http.Handler {
	logger := log.WithContext(ctx)

	if err := view.Register(
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
	); err != nil {
		logger.Fatal(err)
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "coffee-shop",
	})
	if err != nil {
		logger.Fatalw("failed to create the prometheus stats exporter", "err", err.Error())
	}
	view.RegisterExporter(exporter)
	view.SetReportingPeriod(1 * time.Second)

	return exporter
}
