package coffees

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoffee(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			"create a new type of coffee",
			testCreateCoffee,
		},
		{
			"tests if a coffee matches a string",
			testCoffeeTypeMatch,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testCreateCoffee(t *testing.T) {
	c := NewCoffee(uuid.Nil, "test", 1.73)
	assert.NotNil(t, c)
	assert.IsType(t, time.Time{}, c.CreatedAt)
}

func testCoffeeTypeMatch(t *testing.T) {
	cc := Cappuccino{}

	assert.True(t, cc.Match("cappuccino"))
	assert.False(t, cc.Match(""))
}
