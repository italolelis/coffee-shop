package coffees

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// mockRepo is an in memory coffees repository. It's mostly used for testing
type mockRepo struct {
	data []*Coffee
}

// NewCoffeeReadWriteRepository creates a new instance of mockRepo
func newMockRepo() *mockRepo {
	return &mockRepo{}
}

// FindOneByID find one coffee by ID
func (r *mockRepo) FindOneByID(ctx context.Context, id uuid.UUID) (*Coffee, error) {
	for _, c := range r.data {
		if c.ID == id {
			return c, nil
		}
	}

	return nil, ErrCoffeeNotFound
}

// FindOneByName find one coffee by name
func (r *mockRepo) FindOneByName(ctx context.Context, name string) (*Coffee, error) {
	for _, c := range r.data {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, ErrCoffeeNotFound
}

// FindAll finds all coffees
func (r *mockRepo) FindAll(ctx context.Context) ([]*Coffee, error) {
	return r.data, nil
}

// Add adds a repository coffee
func (r *mockRepo) Add(ctx context.Context, c *Coffee) error {
	r.data = append(r.data, c)
	return nil
}

// Remove removes a repository rule
func (r *mockRepo) Remove(ctx context.Context, id uuid.UUID) error {
	for i, c := range r.data {
		if c.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}

	return ErrCoffeeNotFound
}

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
	r := newMockRepo()
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
	r := newMockRepo()
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
	r := newMockRepo()
	s := NewService(r, r)

	_, err := s.RequestCoffee(ctx, "wrong")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidID, err)
}

func testRequestInexistentCoffee(t *testing.T) {
	ctx := context.Background()
	r := newMockRepo()
	s := NewService(r, r)

	_, err := s.RequestCoffee(ctx, uuid.NewV4().String())
	assert.Error(t, err)
	assert.Equal(t, ErrCoffeeNotFound, err)
}
