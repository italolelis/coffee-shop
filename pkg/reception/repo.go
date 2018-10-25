package reception

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrOrderNotFound is used when the order is not found
	ErrOrderNotFound = errors.New("could not find the order")
)

// WriteRepository represents the write operations for an order
type WriteRepository interface {
	Add(context.Context, *Order) error
	Remove(context.Context, uuid.UUID) error
}

// PostgresWriteRepository is the PostgresSQL repository
type PostgresWriteRepository struct {
	db *sqlx.DB
}

// NewPostgresWriteRepository creates a new instance of PostgresWriteRepository
func NewPostgresWriteRepository(db *sqlx.DB) *PostgresWriteRepository {
	return &PostgresWriteRepository{db}
}

// Add adds a repository rule
func (r *PostgresWriteRepository) Add(ctx context.Context, o *Order) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO orders 
		VALUES (:id, :items, :created_at)
			ON CONFLICT (id) DO
		UPDATE SET (items) = (:items)
	`, o); err != nil {
		return fmt.Errorf("could not insert an order: %s", err.Error())
	}

	return nil
}

// Remove removes a repository rule
func (r *PostgresWriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id.String()); err != nil {
		return fmt.Errorf("could not delete the order: %s", err.Error())
	}

	return nil
}

// ReadRepository represents the read operations for an order
type ReadRepository interface {
	FindOneByID(context.Context, uuid.UUID) (*Order, error)
}

// PostgresReadRepository is the PostgresSQL repository
type PostgresReadRepository struct {
	db *sqlx.DB
}

// NewPostgresReadRepository creates a new instance of PostgresReadRepository
func NewPostgresReadRepository(db *sqlx.DB) *PostgresReadRepository {
	return &PostgresReadRepository{db}
}

// FindOneByID find one order by ID
func (r *PostgresReadRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var order Order

	if err := r.db.Get(&order, "SELECT * FROM orders WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf("could not find the order %s", err.Error())
	}

	return &order, nil
}
