package inmem

import (
	"context"

	"github.com/italolelis/coffee-shop/internal/app/coffees"
	"github.com/satori/go.uuid"
)

// CoffeeReadWriteRepository is an in memory coffees repository. It's mostly used for testing
type CoffeeReadWriteRepository struct {
	data []*coffees.Coffee
}

// NewCoffeeReadWriteRepository creates a new instance of CoffeeReadWriteRepository
func NewCoffeeReadWriteRepository() *CoffeeReadWriteRepository {
	return &CoffeeReadWriteRepository{}
}

// FindOneByID find one coffee by ID
func (r *CoffeeReadWriteRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*coffees.Coffee, error) {
	for _, c := range r.data {
		if c.ID == id {
			return c, nil
		}
	}

	return nil, coffees.ErrCoffeeNotFound
}

// FindOneByName find one coffee by name
func (r *CoffeeReadWriteRepository) FindOneByName(ctx context.Context, name string) (*coffees.Coffee, error) {
	for _, c := range r.data {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, coffees.ErrCoffeeNotFound
}

func (r *CoffeeReadWriteRepository) FindAll(ctx context.Context) ([]*coffees.Coffee, error) {
	return r.data, nil
}

// Add adds a repository coffee
func (r *CoffeeReadWriteRepository) Add(ctx context.Context, c *coffees.Coffee) error {
	r.data = append(r.data, c)
	return nil
}

// Remove removes a repository rule
func (r *CoffeeReadWriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	for i, c := range r.data {
		if c.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}

	return coffees.ErrCoffeeNotFound
}
