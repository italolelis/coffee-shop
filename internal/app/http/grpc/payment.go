package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/italolelis/coffee-shop/internal/app/payment"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"github.com/italolelis/coffee-shop/internal/pkg/pb"
)

type PaymentHandler struct{}

func (h *PaymentHandler) Pay(ctx context.Context, r *pb.PaymentRequest) (*pb.PaymentConfirmation, error) {
	logger := log.WithContext(ctx).
		Named("payments").
		With("action", "pay").
		With("order_id", r.OrderID)

	m, err := payment.NewMethodFactory(r.Method)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.Parse(r.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order id: %w", err)
	}

	logger.Debug("processing payment")
	c, err := m.Process(payment.OrderRequest{OrderID: orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// save confirmation in storage

	logger.Debug("payment processed")
	return &pb.PaymentConfirmation{ID: c.ID.String(), OrderID: c.OrderID.String()}, nil
}
