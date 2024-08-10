package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSort(t *testing.T) {
	t.Run("Should return new sort", func(t *testing.T) {
		expected := Sort{Direction: ASC, Field: "name"}

		result := NewSort(expected.Direction, expected.Field)

		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
	})
}

func TestSortDirection_IsValid(t *testing.T) {
	t.Run("Should return true for valid sort direction", func(t *testing.T) {
		validDirection := ASC

		result := validDirection.IsValid()

		assert.True(t, result)
	})

	t.Run("Should return false for invalid sort direction", func(t *testing.T) {
		invalidDirection := SortDirection("INVALID")

		result := invalidDirection.IsValid()

		assert.False(t, result)
	})
}
