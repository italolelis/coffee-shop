package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/italolelis/coffee-shop/internal/app/order"
	"github.com/italolelis/coffee-shop/internal/app/storage/inmem"
	"github.com/italolelis/coffee-shop/internal/pkg/pb"
	"github.com/italolelis/coffee-shop/internal/pkg/tracing"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Server represents a REST server
type Server struct {
	s  *http.Server
	oh *OrderHandler
}

// NewServer creates a new Server
func NewServer(cfg Config, pc pb.PaymentClient) *Server {
	orw := inmem.NewOrderReadWrite()
	os := order.NewService(orw, orw, pc)

	return &Server{
		s: &http.Server{
			Addr:         cfg.Addr,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		oh: &OrderHandler{srv: os},
	}
}

// ListenAndServe opens a port on given address
// and listens for GRPC connections
func (s *Server) ListenAndServe(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(tracing.Tracing)
	r.Route("/orders", func(r chi.Router) {
		r.Post("/checkout", http.HandlerFunc(s.oh.Checkout))
		r.Post("/", http.HandlerFunc(s.oh.AddToOrder))
		r.Get("/{orderID}", http.HandlerFunc(s.oh.GetOrder))
	})

	s.s.Handler = r
	s.s.BaseContext = func(l net.Listener) context.Context {
		return ctx
	}

	return s.s.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.s.Shutdown(ctx); err != nil {
		if err = s.s.Close(); err != nil {
			return fmt.Errorf("failed to close probe server gracefully: %w", err)
		}
	}

	return nil
}
