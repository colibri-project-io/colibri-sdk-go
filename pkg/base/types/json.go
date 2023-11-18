package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	ErrorInvalidValue = errors.New("invalid []byte value")
)

type JsonB map[string]interface{}

func (t *JsonB) Scan(value interface{}) error {
	result, valid := value.([]byte)
	if !valid {
		return ErrorInvalidValue
	}

	return json.Unmarshal(result, &t)
}

func (t *JsonB) Value() (driver.Value, error) {
	return json.Marshal(t)
}
