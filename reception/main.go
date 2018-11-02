package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/italolelis/reception/pkg/config"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/reception/pkg/reception"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	db, err := sqlx.Connect("postgres", cfg.Database.DSN)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()

	wRepo := reception.NewPostgresWriteRepository(db)
	rRepo := reception.NewPostgresReadRepository(db)

	// creates the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", reception.CreateOrderHandler(wRepo))
		r.Get("/{id}", reception.GetOrderHandler(rRepo))
	})

	logger.Infow("service running", "port", cfg.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), chi.ServerBaseContext(ctx, r)))
}
