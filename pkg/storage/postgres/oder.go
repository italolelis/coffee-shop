package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/italolelis/coffee-shop/pkg/order"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

// OrderWriteRepository is the PostgresSQL repository
type OrderWriteRepository struct {
	db *sqlx.DB
}

// NewPostgresOrderWriteRepository creates a new instance of OrderWriteRepository
func NewPostgresOrderWriteRepository(db *sqlx.DB) *OrderWriteRepository {
	return &OrderWriteRepository{db}
}

// Add adds a repository rule
func (r *OrderWriteRepository) Add(ctx context.Context, o *order.Order) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO orders 
		VALUES (:id, :items, :created_at, :customer_name)
			ON CONFLICT (id) DO
		UPDATE SET (items, customer_name) = (:items, :customer_name)
	`, o); err != nil {
		return fmt.Errorf("could not insert an order: %s", err.Error())
	}

	return nil
}

// Remove removes a repository rule
func (r *OrderWriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id.String()); err != nil {
		return fmt.Errorf("could not delete the order: %s", err.Error())
	}

	return nil
}

// OrderReadRepository is the PostgresSQL repository
type OrderReadRepository struct {
	db *sqlx.DB
}

// NewPostgresOrderReadRepository creates a new instance of OrderReadRepository
func NewPostgresOrderReadRepository(db *sqlx.DB) *OrderReadRepository {
	return &OrderReadRepository{db}
}

// FindOneByID find one order by ID
func (r *OrderReadRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	var o order.Order

	if err := r.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, order.ErrOrderNotFound
		}

		return nil, fmt.Errorf("could not find the order %s", err.Error())
	}

	return &o, nil
}
