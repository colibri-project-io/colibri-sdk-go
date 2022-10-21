package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSort(t *testing.T) {
	t.Run("Should return a new sort", func(t *testing.T) {
		result := NewSort(ASC, "name")
		assert.NotNil(t, result)
		assert.Equal(t, ASC, result.Direction)
		assert.Equal(t, "name", result.Field)
	})

	t.Run("Should return a new page request", func(t *testing.T) {
		page := uint16(1)
		size := uint16(1)
		sort1 := NewSort(ASC, "name")
		sort2 := NewSort(DESC, "age")

		result := NewPageRequest(page, size, []Sort{sort1, sort2})
		assert.NotNil(t, result)
		assert.Equal(t, page, result.Page)
		assert.Equal(t, size, result.Size)
		assert.Len(t, result.Order, 2)
		assert.Equal(t, sort1, result.Order[0])
		assert.Equal(t, sort2, result.Order[1])
		assert.Equal(t, "name ASC, age DESC", result.GetOrder())
	})
}
