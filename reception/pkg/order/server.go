package order

import (
	"github.com/go-chi/chi"
	"github.com/italolelis/reception/pkg/coffees"
	"github.com/jmoiron/sqlx"
	"github.com/rafaeljesus/rabbus"
)

func NewServer(db *sqlx.DB, es *rabbus.Rabbus) chi.Router {
	wRepo := NewPostgresWriteRepository(db)
	rRepo := NewPostgresReadRepository(db)
	coffeeReadRepo := coffees.NewPostgresReadRepository(db)
	orderHandler := NewHandler(wRepo, rRepo, coffeeReadRepo, es)

	r := chi.NewRouter()
	r.Post("/", orderHandler.CreateOrder)
	r.Get("/{id}", orderHandler.GetOrder)

	return r
}
