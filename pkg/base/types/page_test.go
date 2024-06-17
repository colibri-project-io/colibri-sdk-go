package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPageRequest(t *testing.T) {
	const (
		page uint16 = 1
		size uint16 = 10
	)

	t.Run("Should return a new PageRequest with nil sort list", func(t *testing.T) {
		expected := &PageRequest{Page: page, Size: size}

		result := NewPageRequest(expected.Page, expected.Size, nil)

		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
		assert.Equal(t, "", result.GetOrder())
	})

	t.Run("Should return a new PageRequest with empty sort list", func(t *testing.T) {
		expected := &PageRequest{Page: page, Size: size, Order: []Sort{}}

		result := NewPageRequest(expected.Page, expected.Size, expected.Order)

		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
		assert.Equal(t, "", result.GetOrder())
	})

	t.Run("Should return a new PageRequest with one populated sort list", func(t *testing.T) {
		expected := &PageRequest{Page: page, Size: size, Order: []Sort{NewSort(ASC, "field1")}}

		result := NewPageRequest(expected.Page, expected.Size, expected.Order)

		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
		assert.EqualValues(t, "field1 ASC", result.GetOrder())
	})

	t.Run("Should return a new PageRequest with many populated sort list", func(t *testing.T) {
		expected := &PageRequest{Page: page, Size: size, Order: []Sort{NewSort(ASC, "field1"), NewSort(DESC, "field2")}}

		result := NewPageRequest(expected.Page, expected.Size, expected.Order)

		assert.NotNil(t, result)
		assert.EqualValues(t, expected, result)
		assert.Equal(t, "field1 ASC, field2 DESC", result.GetOrder())
	})
}
