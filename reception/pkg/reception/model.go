package reception

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Order represents a coffee order
type Order struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Items     Items     `json:"items" db:"items"`
}

// Items are a collection of order items
type Items []*Item

// Item represents the order item
type Item struct {
	Type  string  `json:"type" db:"type"`
	Price float32 `json:"price" db:"price"`
}

// NewOrder creates a new instance of Order
func NewOrder() *Order {
	return &Order{
		ID:        uuid.NewV4(),
		CreatedAt: time.Now(),
	}
}

// Value return a driver.Value representation of the order items
func (p Items) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan scans a database json representation into a []Item
func (p *Items) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &p)
	}
	return fmt.Errorf("could not not decode type %T -> %T", src, p)
}
