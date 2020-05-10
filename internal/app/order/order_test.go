package order

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder_AddItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		customerName    string
		items           Items
		errorExpected   bool
		expectedItemQty int
		expectedTotal   float64
	}{
		{
			name:            "add single item empty name",
			customerName:    "test",
			errorExpected:   true,
			expectedItemQty: 0,
			expectedTotal:   0,
			items: Items{
				{
					Name:        "",
					Qty:         2,
					ServingSize: "L",
					Price:       2.60,
				},
			},
		},
		{
			name:            "add single item empty serving size",
			customerName:    "test",
			errorExpected:   true,
			expectedItemQty: 0,
			expectedTotal:   0,
			items: Items{
				{
					Name:        "latte",
					Qty:         2,
					ServingSize: "",
					Price:       2.60,
				},
			},
		},
		{
			name:            "add two same items",
			customerName:    "test",
			errorExpected:   false,
			expectedItemQty: 1,
			expectedTotal:   7.80,
			items: Items{
				{
					Name:        "cappuccino",
					Qty:         2,
					ServingSize: "L",
					Price:       2.60,
				},
				{
					Name:        "cappuccino",
					Qty:         1,
					ServingSize: "L",
					Price:       2.60,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := New(tt.customerName)

			err := o.AddItems(tt.items)
			if !tt.errorExpected {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedItemQty, len(o.Items))
			assert.Equal(t, tt.expectedTotal, math.Floor(o.Total()*100)/100)
		})
	}
}
