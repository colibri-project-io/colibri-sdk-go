package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"
)

// NullBoll for empty boolean field
type NullBool sql.NullBool

// NullFloat64 for empty float field
type NullFloat64 sql.NullFloat64

// NullInt16 for empty int16 field
type NullInt16 sql.NullInt16

// NullInt32 for empty int32 field
type NullInt32 sql.NullInt32

// NullInt64 for empty int64 field
type NullInt64 sql.NullInt64

// NullString for empty string field
type NullString sql.NullString

// NullTime for empty date/time field
type NullTime sql.NullTime

func (t *NullBool) Scan(value interface{}) error {
	var i sql.NullBool
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullBool{i.Bool, false}
	} else {
		*t = NullBool{i.Bool, true}
	}

	return nil
}

func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Bool, nil
}

func (t NullBool) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Bool)
}

func (t *NullBool) UnmarshalJSON(data []byte) error {
	var ptr *bool
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Bool = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullFloat64) Scan(value interface{}) error {
	var i sql.NullFloat64
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullFloat64{i.Float64, false}
	} else {
		*t = NullFloat64{i.Float64, true}
	}

	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Float64, nil
}

func (t NullFloat64) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Float64)
}

func (t *NullFloat64) UnmarshalJSON(data []byte) error {
	var ptr *float64
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Float64 = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullInt16) Scan(value interface{}) error {
	var i sql.NullInt16
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullInt16{i.Int16, false}
	} else {
		*t = NullInt16{i.Int16, true}
	}

	return nil
}

func (n NullInt16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Int16, nil
}

func (t NullInt16) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Int16)
}

func (t *NullInt16) UnmarshalJSON(data []byte) error {
	var ptr *int16
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Int16 = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullInt32) Scan(value interface{}) error {
	var i sql.NullInt32
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullInt32{i.Int32, false}
	} else {
		*t = NullInt32{i.Int32, true}
	}

	return nil
}

func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Int32, nil
}

func (t NullInt32) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Int32)
}

func (t *NullInt32) UnmarshalJSON(data []byte) error {
	var ptr *int32
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Int32 = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullInt64{i.Int64, false}
	} else {
		*t = NullInt64{i.Int64, true}
	}

	return nil
}

func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Int64, nil
}

func (t NullInt64) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Int64)
}

func (t *NullInt64) UnmarshalJSON(data []byte) error {
	var ptr *int64
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Int64 = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullString) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullString{i.String, false}
	} else {
		*t = NullString{i.String, true}
	}

	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}

	return ns.String, nil
}

func (t NullString) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.String)
}

func (t *NullString) UnmarshalJSON(data []byte) error {
	var ptr *string
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.String = *ptr
	} else {
		t.Valid = false
	}

	return nil
}

func (t *NullTime) Scan(value interface{}) error {
	var i sql.NullTime
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullTime{i.Time, false}
	} else {
		*t = NullTime{i.Time, true}
	}

	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Time, nil
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(t.Time)
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	var ptr *time.Time
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr != nil {
		t.Valid = true
		t.Time = *ptr
	} else {
		t.Valid = false
	}

	return nil
}
