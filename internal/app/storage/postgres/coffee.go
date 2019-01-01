package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/italolelis/coffee-shop/internal/app/coffees"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

// CoffeeWriteRepository is the coffee PostgresSQL repository
type CoffeeWriteRepository struct {
	db *sqlx.DB
}

// NewCoffeeWriteRepository creates a new instance of OrderWriteRepository
func NewCoffeeWriteRepository(db *sqlx.DB) *CoffeeWriteRepository {
	return &CoffeeWriteRepository{db}
}

// Add adds a repository rule
func (r *CoffeeWriteRepository) Add(ctx context.Context, o *coffees.Coffee) error {
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
func (r *CoffeeWriteRepository) Remove(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM coffees WHERE id = $1", id.String()); err != nil {
		return fmt.Errorf("could not delete the coffee: %s", err.Error())
	}

	return nil
}

// CoffeeReadRepository is the coffee PostgresSQL repository
type CoffeeReadRepository struct {
	db *sqlx.DB
}

// NewCoffeeReadRepository creates a new instance of OrderReadRepository
func NewCoffeeReadRepository(db *sqlx.DB) *CoffeeReadRepository {
	return &CoffeeReadRepository{db}
}

// FindOneByID find one order by ID
func (r *CoffeeReadRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*coffees.Coffee, error) {
	var c coffees.Coffee

	if err := r.db.GetContext(ctx, &c, "SELECT * FROM coffees WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, coffees.ErrCoffeeNotFound
		}

		return nil, fmt.Errorf("could not find the coffee %s", err.Error())
	}

	return &c, nil
}

// FindOneByName find one order by ID
func (r *CoffeeReadRepository) FindOneByName(ctx context.Context, name string) (*coffees.Coffee, error) {
	var c coffees.Coffee

	if err := r.db.GetContext(ctx, &c, "SELECT * FROM coffees WHERE name = $1 LIMIT 1", name); err != nil {
		if err == sql.ErrNoRows {
			return nil, coffees.ErrCoffeeNotFound
		}

		return nil, fmt.Errorf("could not find the coffee %s", err.Error())
	}

	return &c, nil
}

// FindAll finds all coffees
func (r *CoffeeReadRepository) FindAll(ctx context.Context) ([]*coffees.Coffee, error) {
	var data []*coffees.Coffee

	if err := r.db.Select(&data, "SELECT * FROM coffees"); err != nil {
		return nil, fmt.Errorf("could not find coffees %s", err.Error())
	}

	return data, nil
}
