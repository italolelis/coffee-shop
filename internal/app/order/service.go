package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/italolelis/coffee-shop/internal/pkg/pb"
	"github.com/italolelis/coffee-shop/internal/pkg/tracing"
)

var ErrNotFound = errors.New("order not found")

type Reader interface {
	FetchByID(context.Context, uuid.UUID) (*Order, error)
}

type Writer interface {
	Add(context.Context, *Order) error
}

type Service interface {
	Checkout(context.Context, CheckoutCommand) (uuid.UUID, error)
	AddToOrder(context.Context, AddToOrderCommand) (uuid.UUID, error)
	Fetch(context.Context, uuid.UUID) (*Order, error)
}

type CheckoutCommand struct {
	CustomerName  string `json:"customer_name"`
	PaymentMethod string `json:"payment_method"`
}

type AddToOrderCommand struct {
	CustomerName string `json:"customer_name"`
	Items        Items  `json:"items"`
}

type ServiceImp struct {
	w  Writer
	r  Reader
	pc pb.PaymentClient
}

func NewService(w Writer, r Reader, pc pb.PaymentClient) *ServiceImp {
	return &ServiceImp{
		w:  w,
		r:  r,
		pc: pc,
	}
}

func (s *ServiceImp) Checkout(ctx context.Context, cmd CheckoutCommand) (uuid.UUID, error) {
	ctx, span := tracing.Start(ctx, "service/order/checkout")
	defer span.End()

	o, err := s.r.FetchByID(ctx, uuid.NewSHA1(uuid.NameSpaceOID, []byte(cmd.CustomerName)))
	if err != nil {
		return uuid.Nil, err
	}

	_, err = s.pc.Pay(ctx, &pb.PaymentRequest{
		Method:  cmd.PaymentMethod,
		OrderID: o.ID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed paying order: %w", err)
	}

	if err := s.w.Add(ctx, o); err != nil {
		return uuid.Nil, fmt.Errorf("failed saving order: %w", err)
	}

	return o.ID, nil
}

func (s *ServiceImp) AddToOrder(ctx context.Context, cmd AddToOrderCommand) (uuid.UUID, error) {
	ctx, span := tracing.Start(ctx, "service/order/add-to-order")
	defer span.End()

	o := New(cmd.CustomerName)
	if err := o.AddItems(cmd.Items); err != nil {
		return uuid.Nil, fmt.Errorf("failed adding items to orders: %w", err)
	}

	if err := s.w.Add(ctx, o); err != nil {
		return uuid.Nil, fmt.Errorf("failed saving order: %w", err)
	}

	return o.ID, nil
}

func (s *ServiceImp) Fetch(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.r.FetchByID(ctx, id)
}
