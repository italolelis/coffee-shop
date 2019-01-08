package coffees

import (
	"context"
	"github.com/italolelis/coffee-shop/internal/app/storage/inmem"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			"create a new coffee successfully",
			testCreateCoffeeSuccessfully,
		},
		{
			"create a new coffee with missing arguments",
			testCoffeeCreatingWithMissingArgs,
		},
		{
			"requests a coffee with an invalid ID",
			testRequestCoffeeWithWrongID,
		},
		{
			"requests a coffee that doesn't exists",
			testRequestInexistentCoffee,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testCreateCoffeeSuccessfully(t *testing.T) {
	ctx := context.Background()
	r := inmem.NewCoffeeReadWriteRepository()
	s := NewService(r, r)

	id, err := s.CreateCoffee(ctx, "test", 1.70)
	assert.NoError(t, err)

	assert.NotNil(t, id)
	assert.IsType(t, uuid.UUID{}, id)

	c, err := s.RequestCoffee(ctx, id.String())
	assert.NoError(t, err)

	assert.Equal(t, "test", c.Name)
}

func testCoffeeCreatingWithMissingArgs(t *testing.T) {
	ctx := context.Background()
	r := inmem.NewCoffeeReadWriteRepository()
	s := NewService(r, r)

	id, err := s.CreateCoffee(ctx, "", 1.70)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidName, err)
	assert.IsType(t, uuid.Nil, id)

	id, err = s.CreateCoffee(ctx, "test", 0)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPrice, err)
	assert.IsType(t, uuid.Nil, id)
}

func testRequestCoffeeWithWrongID(t *testing.T) {
	ctx := context.Background()
	r := inmem.NewCoffeeReadWriteRepository()
	s := NewService(r, r)

	_, err := s.RequestCoffee(ctx, "wrong")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidID, err)
}

func testRequestInexistentCoffee(t *testing.T) {
	ctx := context.Background()
	r := inmem.NewCoffeeReadWriteRepository()
	s := NewService(r, r)

	_, err := s.RequestCoffee(ctx, uuid.NewV4().String())
	assert.Error(t, err)
	assert.Equal(t, ErrCoffeeNotFound, err)
}
