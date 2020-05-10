package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/italolelis/coffee-shop/internal/pkg/pb"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/grpctrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	Addr string
}

// Server represents a GRPC server
type Server struct {
	cfg Config
	ph  *PaymentHandler
	g   *grpc.Server
}

// NewServer creates a new Server
func NewServer(cfg Config, tp trace.Tracer) *Server {
	return &Server{
		cfg: cfg,
		g: grpc.NewServer(
			grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor(tp)),
			grpc.StreamInterceptor(grpctrace.StreamServerInterceptor(tp)),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Timeout: 30 * time.Second,
			}),
		),
		ph: &PaymentHandler{},
	}
}

// ListenAndServe opens a port on given address
// and listens for GRPC connections
func (s *Server) ListenAndServe(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on address: %w", err)
	}

	pb.RegisterPaymentServer(s.g, s.ph)

	return s.g.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) {
	s.g.GracefulStop()
}
