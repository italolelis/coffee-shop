package coffees

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func NewServer(db *sqlx.DB) chi.Router {
	wRepo := NewPostgresWriteRepository(db)
	rRepo := NewPostgresReadRepository(db)
	handler := NewHandler(wRepo, rRepo)

	r := chi.NewRouter()
	r.Post("/", handler.CreateCoffee)
	r.Get("/{id}", handler.GetCoffee)

	return r
}
