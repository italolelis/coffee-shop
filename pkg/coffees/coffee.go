package coffees

import (
	"context"
	"time"

	"github.com/satori/go.uuid"
)

// Coffee represents a coffee type
type Coffee struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	Price     float32   `json:"price" db:"price"`
}

// NewCoffee creates a new instance of Coffee
func NewCoffee(id uuid.UUID, name string, price float32) *Coffee {
	return &Coffee{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      name,
		Price:     price,
	}
}

func NextCoffeeID() uuid.UUID {
	return uuid.NewV4()
}

type CoffeeType interface {
	Brew(context.Context)
	Match(string) bool
}

type Espresso struct{}

func (e *Espresso) Brew(ctx context.Context) {
	time.Sleep(3 * time.Second)
}

func (e *Espresso) Match(coffeeType string) bool {
	return coffeeType == "espresso"
}

type Cappuccino struct{}

func (c *Cappuccino) Brew(ctx context.Context) {
	time.Sleep(6 * time.Second)
}

func (c *Cappuccino) Match(coffeeType string) bool {
	return coffeeType == "cappuccino"
}

type Latte struct{}

func (l *Latte) Brew(ctx context.Context) {
	time.Sleep(8 * time.Second)
}

func (l *Latte) Match(coffeeType string) bool {
	return coffeeType == "latte"
}
