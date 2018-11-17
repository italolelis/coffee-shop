package coffees

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Coffee represents a coffee type
type Coffee struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	Price     float32   `json:"price" db:"price"`
}

// NewCoffee creates a new instance of Coffee
func NewCoffee() *Coffee {
	return &Coffee{
		ID:        uuid.NewV4(),
		CreatedAt: time.Now(),
	}
}
