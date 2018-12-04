package coffees

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

var (
	// ErrCoffeeNotFound is used when the coffee is not found
	ErrCoffeeNotFound = errors.New("could not find the coffee")
)

// WriteRepository represents the write operations for a coffee
type WriteRepository interface {
	Add(context.Context, *Coffee) error
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
func (r *PostgresWriteRepository) Add(ctx context.Context, o *Coffee) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO coffees 
		VALUES (:id, :name, :price, :created_at)
			ON CONFLICT (id) DO
		UPDATE SET (name, price) = (:name, :price)
	`, o); err != nil {
		return fmt.Errorf("could not insert an order: %s", err.Error())
	}

	return nil
}

// Remove removes a repository rule
func (r *PostgresWriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM coffees WHERE id = $1", id.String()); err != nil {
		return fmt.Errorf("could not delete the coffee: %s", err.Error())
	}

	return nil
}

// ReadRepository represents the read operations for a coffee
type ReadRepository interface {
	FindOneByID(context.Context, uuid.UUID) (*Coffee, error)
	FindOneByName(context.Context, string) (*Coffee, error)
	FindAll(context.Context) ([]*Coffee, error)
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
func (r *PostgresReadRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*Coffee, error) {
	var c Coffee

	if err := r.db.GetContext(ctx, &c, "SELECT * FROM coffees WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCoffeeNotFound
		}

		return nil, fmt.Errorf("could not find the coffee %s", err.Error())
	}

	return &c, nil
}

// FindOneByID find one order by ID
func (r *PostgresReadRepository) FindOneByName(ctx context.Context, name string) (*Coffee, error) {
	var c Coffee

	if err := r.db.GetContext(ctx, &c, "SELECT * FROM coffees WHERE name = $1 LIMIT 1", name); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCoffeeNotFound
		}

		return nil, fmt.Errorf("could not find the coffee %s", err.Error())
	}

	return &c, nil
}

// FindAll finds all coffees
func (r *PostgresReadRepository) FindAll(ctx context.Context) ([]*Coffee, error) {
	var data []*Coffee

	if err := r.db.Select(&data, "SELECT * FROM coffees"); err != nil {
		return nil, fmt.Errorf("could not find coffees %s", err.Error())
	}

	return data, nil
}
