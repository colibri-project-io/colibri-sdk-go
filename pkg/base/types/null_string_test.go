package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullString(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullString
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, "", result.String)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := "string test"

		var result NullString
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.String)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullString{"string test", true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.String, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullString{"", false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullString{"", false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullString{"string test", true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"string test\"", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullString
		err := result.UnmarshalJSON([]byte("\"string test\""))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, "string test", result.String)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullString
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, "", result.String)
	})
}
