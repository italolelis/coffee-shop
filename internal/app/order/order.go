package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		ID           uuid.UUID `json:"id" db:"id"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		CustomerName string    `json:"customer" db:"customer"`
		Items        Items     `json:"items" db:"items"`
	}

	Items []*Item

	Item struct {
		Name        string  `json:"name" db:"name"`
		ServingSize string  `json:"serving_size" db:"serving_size"`
		Price       float64 `json:"price" db:"price"`
		Qty         int     `json:"qty" db:"qty"`
	}
)

func New(customerName string) *Order {
	return &Order{
		ID:           uuid.NewSHA1(uuid.NameSpaceOID, []byte(customerName)),
		CreatedAt:    time.Now().UTC(),
		CustomerName: customerName,
		Items:        make([]*Item, 0),
	}
}

func (o *Order) AddItems(items Items) error {
	for _, i := range items {
		if err := o.AddItem(i); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) AddItem(i *Item) error {
	if i.Name == "" {
		return errors.New("item name can't be empty")
	}

	if i.ServingSize == "" {
		return errors.New("serving size can't be empty")
	}

	for _, existingItem := range o.Items {
		if existingItem.Name == i.Name && existingItem.ServingSize == i.ServingSize {
			existingItem.Qty += i.Qty
			return nil
		}
	}

	o.Items = append(o.Items, i)

	return nil
}

func (o *Order) Total() float64 {
	var total float64
	for _, i := range o.Items {
		total += i.Price * float64(i.Qty)
	}

	return total
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
