package inmem

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/italolelis/coffee-shop/internal/app/order"
	"github.com/italolelis/coffee-shop/internal/pkg/tracing"
)

type OrderReadWrite struct {
	mux    *sync.RWMutex
	orders map[uuid.UUID]*order.Order
}

func NewOrderReadWrite() *OrderReadWrite {
	return &OrderReadWrite{mux: &sync.RWMutex{}, orders: make(map[uuid.UUID]*order.Order, 0)}
}

func (r *OrderReadWrite) FetchByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	ctx, span := tracing.Start(ctx, "storage/order/fetch-by-id")
	defer span.End()

	o, ok := r.orders[id]
	if !ok {
		return nil, order.ErrNotFound
	}

	return o, nil
}

func (r *OrderReadWrite) Add(ctx context.Context, o *order.Order) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	ctx, span := tracing.Start(ctx, "storage/order/add")
	defer span.End()

	existingOrder, ok := r.orders[o.ID]
	if !ok {
		r.orders[o.ID] = o
		return nil
	}

	existingOrder.AddItems(o.Items)

	return nil
}
