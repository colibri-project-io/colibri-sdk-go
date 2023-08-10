package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsoTime(t *testing.T) {
	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := time.Date(1, time.January, 1, 10, 20, 30, 0, time.UTC)

		result, err := IsoTime(expected).Value()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := IsoTime(time.Date(0, time.January, 1, 10, 20, 30, 0, time.UTC))

		json, err := expected.MarshalJSON()
		result := string(json)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"10:20:30\"", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result IsoTime
		err := result.UnmarshalJSON([]byte("\"10:20:30\""))

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, IsoTime(time.Date(0, time.January, 1, 10, 20, 30, 0, time.UTC)), result)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result IsoTime
		err := result.UnmarshalJSON([]byte("invalid"))

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, IsoTime(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result)
	})
}
