package payment

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderRequest struct {
	OrderID uuid.UUID
	Total   float64
}

type Method interface {
	Process(OrderRequest) (*Confirmation, error)
}

type Confirmation struct {
	ID      uuid.UUID `json:"id" db:"id"`
	OrderID uuid.UUID `json:"order_id" db:"order_id"`
	PayedAt time.Time `json:"payed_at" db:"payed_at"`
}

func NewConfirmation(orderID uuid.UUID) *Confirmation {
	return &Confirmation{ID: uuid.New(), OrderID: orderID, PayedAt: time.Now()}
}

type MethodFactory struct{}

func NewMethodFactory(method string) (Method, error) {
	switch method {
	case "credit_card":
		return &CreditCard{}, nil
	case "apple_pay":
		return &ApplePay{}, nil
	default:
		return nil, errors.New("payment method not supported")
	}
}

type CreditCard struct{}

func (c *CreditCard) Process(o OrderRequest) (*Confirmation, error) {
	// pretend to connect to some credit card provider
	return NewConfirmation(o.OrderID), nil
}

type ApplePay struct{}

func (a *ApplePay) Process(o OrderRequest) (*Confirmation, error) {
	// pretend to connect to apple pay
	return NewConfirmation(o.OrderID), nil
}
