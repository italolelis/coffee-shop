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
	"github.com/italolelis/kit/metric"
	"github.com/italolelis/kit/stream"
	"github.com/italolelis/kit/trace"
	"github.com/italolelis/reception/pkg/coffees"
	"github.com/italolelis/reception/pkg/config"
	"github.com/italolelis/reception/pkg/order"
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

	// loads the configuration from the environment
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	log.SetLevel(cfg.LogLevel)

	db, dbClose := setupDatabase(ctx, cfg.Database)
	defer dbClose()

	metricsHandler := metric.Setup(ctx, cfg.Metrics)

	eventStream, streamClose := stream.Setup(ctx, cfg.EventStream)
	defer streamClose()

	flush := trace.Setup(ctx, cfg.Tracing)
	defer flush()

	// creates the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/metrics", metricsHandler)
	r.Mount("/orders", order.NewServer(db, eventStream))
	r.Mount("/coffees", coffees.NewServer(db))

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
