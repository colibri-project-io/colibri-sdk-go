package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullBool(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullBool
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, false, result.Bool)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullBool
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, false, result.Bool)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := true

		var result NullBool
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Bool)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullBool{true, true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Bool, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullBool{false, false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullBool{false, false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullBool{true, true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "true", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullBool
		err := result.UnmarshalJSON([]byte("true"))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, true, result.Bool)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullBool
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, false, result.Bool)
	})
}

func TestNullFloat64(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullFloat64
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, 0.00, result.Float64)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullFloat64
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, 0.00, result.Float64)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := 123.45

		var result NullFloat64
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Float64)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullFloat64{123.45, true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Float64, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullFloat64{0.00, false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullFloat64{0.00, false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullFloat64{123.45, true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "123.45", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullFloat64
		err := result.UnmarshalJSON([]byte("123.45"))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, 123.45, result.Float64)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullFloat64
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, 0.00, result.Float64)
	})
}

func TestNullInt16(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullInt16
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int16(0), result.Int16)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullInt16
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int16(0), result.Int16)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := int16(123)

		var result NullInt16
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Int16)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullInt16{12, true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Int16, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullInt16{0, false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullInt16{0, false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullInt16{123, true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "123", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullInt16
		err := result.UnmarshalJSON([]byte("123"))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, int16(123), result.Int16)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullInt16
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int16(0), result.Int16)
	})
}

func TestNullInt32(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullInt32
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullInt32
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := int32(123)

		var result NullInt32
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Int32)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullInt32{12, true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Int32, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullInt32{0, false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullInt32{0, false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullInt32{123, true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "123", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullInt32
		err := result.UnmarshalJSON([]byte("123"))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, int32(123), result.Int32)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullInt32
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int32(0), result.Int32)
	})
}

func TestNullInt64(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullInt64
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int64(0), result.Int64)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullInt64
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int64(0), result.Int64)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := int64(123)

		var result NullInt64
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Int64)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullInt64{123, true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Int64, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullInt64{0, false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullInt64{0, false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullInt64{123, true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "123", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullInt64
		err := result.UnmarshalJSON([]byte("123"))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, int64(123), result.Int64)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullInt64
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, int64(0), result.Int64)
	})
}

func TestNullIsoDate(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullIsoDate
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		var result NullIsoDate
		err := result.Scan("invalid")

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		expected := time.Now()

		var result NullIsoDate
		err := result.Scan(expected)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, expected, result.Time)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullIsoDate{time.Now(), true}

		result, err := expected.Value()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Time, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullIsoDate{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullIsoDate{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), false}

		json, err := expected.MarshalJSON()
		result := string(json)

		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullIsoDate{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), true}

		json, err := expected.MarshalJSON()
		result := string(json)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"2022-01-01\"", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullIsoDate
		err := result.UnmarshalJSON([]byte("\"2022-01-01\""))

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullIsoDate
		err := result.UnmarshalJSON([]byte("invalid"))

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})
}

func TestNullIsoTime(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullIsoTime
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		var result NullIsoTime
		err := result.Scan("invalid")

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		expected := time.Now()

		var result NullIsoTime
		err := result.Scan(expected)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, expected, result.Time)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullIsoTime{time.Now(), true}

		result, err := expected.Value()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Time, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullIsoTime{time.Time(time.Date(1, time.January, 1, 10, 20, 30, 0, time.UTC)), false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullIsoTime{time.Time(time.Date(1, time.January, 1, 10, 20, 30, 0, time.UTC)), false}

		json, err := expected.MarshalJSON()
		result := string(json)

		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullIsoTime{time.Time(time.Date(0, time.January, 1, 10, 20, 30, 0, time.UTC)), true}

		json, err := expected.MarshalJSON()
		result := string(json)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"10:20:30\"", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullIsoTime
		err := result.UnmarshalJSON([]byte("\"10:20:30\""))

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, time.Time(time.Date(0, time.January, 1, 10, 20, 30, 0, time.UTC)), result.Time)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullIsoTime
		err := result.UnmarshalJSON([]byte("invalid"))

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})
}

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

func TestNullDateTime(t *testing.T) {
	t.Run("Should error when scan with a nil value", func(t *testing.T) {
		var result NullDateTime
		err := result.Scan(nil)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should error when scan with a invalid value", func(t *testing.T) {
		value := "invalid"

		var result NullDateTime
		err := result.Scan(value)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should scan with a valid value", func(t *testing.T) {
		value := time.Now()

		var result NullDateTime
		err := result.Scan(value)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, value, result.Time)
	})

	t.Run("Should get value with a valid value", func(t *testing.T) {
		expected := NullDateTime{time.Now(), true}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.Time, result)
	})

	t.Run("Should return nil when get value with a invalid value", func(t *testing.T) {
		expected := NullDateTime{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), false}

		result, err := expected.Value()
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return null when get json value with a invalid value", func(t *testing.T) {
		expected := NullDateTime{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), false}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.Equal(t, "null", result)
	})

	t.Run("Should get json value with a valid value", func(t *testing.T) {
		expected := NullDateTime{time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), true}

		json, err := expected.MarshalJSON()
		result := string(json)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "\"2022-01-01T00:00:00Z\"", result)
	})

	t.Run("Should get value with a valid json", func(t *testing.T) {
		var result NullDateTime
		err := result.UnmarshalJSON([]byte("\"2022-01-01T00:00:00Z\""))
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, true, result.Valid)
		assert.Equal(t, time.Time(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})

	t.Run("Should return error when get value with a invalid json", func(t *testing.T) {
		var result NullDateTime
		err := result.UnmarshalJSON([]byte("invalid"))
		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, false, result.Valid)
		assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), result.Time)
	})
}
