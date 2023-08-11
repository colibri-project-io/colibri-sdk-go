package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// IsoTime struct
type IsoTime time.Time

// Value converts iso time to sql driver value
func (t IsoTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// MarshalJSON converts iso time to json string format
func (t IsoTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.TimeOnly))
}

// UnmarshalJSON converts json string to iso time
func (t *IsoTime) UnmarshalJSON(data []byte) error {
	var ptr *string
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr == nil {
		return nil
	}

	parsedTime, err := time.Parse(time.TimeOnly, *ptr)
	if err != nil {
		return err
	}

	*t = IsoTime(parsedTime)
	return nil
}
