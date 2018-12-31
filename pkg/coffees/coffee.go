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

// NextCoffeeID generates the next coffee ID
func NextCoffeeID() uuid.UUID {
	return uuid.NewV4()
}

// CoffeeType is an interface for the different coffee preparations
type CoffeeType interface {
	Brew(context.Context)
	Match(string) bool
}

// Espresso is a type of coffee
type Espresso struct{}

// Brew the coffee
func (e *Espresso) Brew(ctx context.Context) {
	time.Sleep(3 * time.Second)
}

// Match checks if the given string is the current coffee
func (e *Espresso) Match(coffeeType string) bool {
	return coffeeType == "espresso"
}

// Cappuccino is a type of coffee
type Cappuccino struct{}

// Brew the coffee
func (c *Cappuccino) Brew(ctx context.Context) {
	time.Sleep(6 * time.Second)
}

// Match checks if the given string is the current coffee
func (c *Cappuccino) Match(coffeeType string) bool {
	return coffeeType == "cappuccino"
}

// Latte is a type of coffee
type Latte struct{}

// Brew the coffee
func (l *Latte) Brew(ctx context.Context) {
	time.Sleep(8 * time.Second)
}

// Match checks if the given string is the current coffee
func (l *Latte) Match(coffeeType string) bool {
	return coffeeType == "latte"
}
