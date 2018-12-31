package coffees

import (
	"context"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
)

var (
	// ErrInvalidName is returned when the provided name is blank.
	ErrInvalidName = errors.New("invalid name provided")

	// ErrInvalidPrice is returned when the provided price is invalid.
	ErrInvalidPrice = errors.New("invalid price provided")

	// ErrInvalidID is returned when the provided ID is invalid.
	ErrInvalidID = errors.New("invalid id provided")

	// ErrCoffeeNotFound is returned when the coffee is not found.
	ErrCoffeeNotFound = errors.New("could not find the coffee")
)

// WriteRepository represents the write operations for a coffee
type WriteRepository interface {
	Add(context.Context, *Coffee) error
	Remove(context.Context, uuid.UUID) error
}

// ReadRepository represents the read operations for a coffee
type ReadRepository interface {
	FindOneByID(context.Context, uuid.UUID) (*Coffee, error)
	FindOneByName(context.Context, string) (*Coffee, error)
	FindAll(context.Context) ([]*Coffee, error)
}

// Service is the interface that provides coffee methods.
type Service interface {
	CreateCoffee(ctx context.Context, name string, price float32) (uuid.UUID, error)
	RequestCoffee(ctx context.Context, coffeeID string) (*Coffee, error)
}

type ServiceImp struct {
	wRepo WriteRepository
	rRepo ReadRepository
}

func NewService(wRepo WriteRepository, rRepo ReadRepository) *ServiceImp {
	return &ServiceImp{wRepo: wRepo, rRepo: rRepo}
}

func (s *ServiceImp) CreateCoffee(ctx context.Context, name string, price float32) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, ErrInvalidName
	}

	if price == 0.0 {
		return uuid.Nil, ErrInvalidPrice
	}

	c := NewCoffee(NextCoffeeID(), name, price)
	if err := s.wRepo.Add(ctx, c); err != nil {
		return uuid.Nil, fmt.Errorf("could not save your coffee: %s", err.Error())
	}

	return c.ID, nil
}

func (s *ServiceImp) RequestCoffee(ctx context.Context, coffeeID string) (*Coffee, error) {
	id, err := uuid.FromString(coffeeID)
	if err != nil {
		return nil, ErrInvalidID
	}

	return s.rRepo.FindOneByID(ctx, id)
}
