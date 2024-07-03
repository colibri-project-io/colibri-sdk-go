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
