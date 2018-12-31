package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/italolelis/coffee-shop/internal/proto/order"
	"github.com/italolelis/coffee-shop/pkg/coffees"
	"github.com/rafaeljesus/rabbus"
	"github.com/satori/go.uuid"
	"go.opencensus.io/trace"
)

const (
	exchangeName = "orders"
	exchangeType = "topic"
)

var (
	// ErrInvalidName is returned when the provided name is blank.
	ErrInvalidName = errors.New("invalid name provided")

	// ErrEmptyOrder is returned when the provided order has no items.
	ErrEmptyOrder = errors.New("the order cannot be empty")

	// ErrInvalidID is returned when the provided ID is invalid.
	ErrInvalidID = errors.New("invalid id provided")

	// ErrOrderNotFound is returned when the order is not found.
	ErrOrderNotFound = errors.New("could not find the order")

	// ErrUnsupportedCoffeeType is returned when the order is not found.
	ErrUnsupportedCoffeeType = errors.New("the requested coffee is not supported yet")
)

// WriteRepository represents the write operations for an order
type WriteRepository interface {
	Add(context.Context, *Order) error
	Remove(context.Context, uuid.UUID) error
}

// ReadRepository represents the read operations for an order
type ReadRepository interface {
	FindOneByID(context.Context, uuid.UUID) (*Order, error)
}

// Service is the interface that provides coffee methods.
type Service interface {
	CreateOrder(ctx context.Context, name string, items Items) (uuid.UUID, error)
	RequestOrder(ctx context.Context, orderID string) (*Order, error)
}

// ServiceImp is the implementation of the service
type ServiceImp struct {
	wRepo   WriteRepository
	rRepo   ReadRepository
	coffees coffees.ReadRepository
	es      *rabbus.Rabbus
}

// NewService creates a new instance of ServiceImp
func NewService(wRepo WriteRepository, rRepo ReadRepository, coffees coffees.ReadRepository, es *rabbus.Rabbus) *ServiceImp {
	return &ServiceImp{wRepo: wRepo, rRepo: rRepo, coffees: coffees, es: es}
}

// CreateOrder creates a new order for the barista
func (s *ServiceImp) CreateOrder(ctx context.Context, name string, items Items) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, ErrInvalidName
	}

	if len(items) <= 0 {
		return uuid.Nil, ErrEmptyOrder
	}

	for _, i := range items {
		_, err := s.coffees.FindOneByName(ctx, i.Type)
		if err != nil {
			return uuid.Nil, ErrUnsupportedCoffeeType
		}
	}

	o := NewOrder(NextOrderID(), name)
	o.AddItems(items)

	if err := s.wRepo.Add(ctx, o); err != nil {
		return uuid.Nil, fmt.Errorf("could not save your coffee: %s", err.Error())
	}

	ev := order.Created{
		ID:           o.ID.String(),
		CustomerName: o.CustomerName,
	}
	for _, i := range o.Items {
		ev.Items = append(ev.Items, &order.OrderItem{Type: i.Type, Size: i.Size})
	}

	if err := s.sendEvent(ctx, "orders.created", &ev); err != nil {
		return uuid.Nil, err
	}

	return o.ID, nil
}

// RequestOrder retrieves an order
func (s *ServiceImp) RequestOrder(ctx context.Context, orderID string) (*Order, error) {
	id, err := uuid.FromString(orderID)
	if err != nil {
		return nil, ErrInvalidID
	}

	return s.rRepo.FindOneByID(ctx, id)
}

func (s *ServiceImp) sendEvent(ctx context.Context, key string, payload proto.Message) error {
	ctx, span := trace.StartSpan(ctx, "rabbitmq:send", trace.WithSpanKind(trace.SpanKindClient))
	span.AddAttributes(trace.StringAttribute("exchange", exchangeName), trace.StringAttribute("key", key))
	defer span.End()

	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal the message before sending to the message stream: %s", err.Error())
	}

	msg := rabbus.Message{
		Exchange: exchangeName,
		Kind:     exchangeType,
		Key:      key,
		Payload:  data,
	}

	s.es.EmitAsync() <- msg

	return nil
}
