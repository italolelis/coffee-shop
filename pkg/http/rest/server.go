package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/italolelis/coffee-shop/pkg/coffees"
	"github.com/italolelis/coffee-shop/pkg/order"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	router chi.Router
}

func NewServer(cs coffees.Service, os order.Service, metricsHandler http.Handler) *Server {
	// creates the router and register the handlers
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	ch := coffeeHandler{cs}
	oh := orderHandler{os}

	r.Handle("/metrics", metricsHandler)
	r.Mount("/orders", oh.router())
	r.Mount("/coffees", ch.router())

	return &Server{router: r}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
