package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// IsoDate struct
type IsoDate time.Time

// Value converts iso date to sql driver value
func (t IsoDate) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// MarshalJSON converts iso date to json string format
func (t IsoDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.DateOnly))
}

// UnmarshalJSON converts json string to iso date
func (t *IsoDate) UnmarshalJSON(data []byte) error {
	var ptr *string
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}

	if ptr == nil {
		return nil
	}

	parsedDate, err := time.Parse(time.DateOnly, *ptr)
	if err != nil {
		return err
	}

	*t = IsoDate(parsedDate)
	return nil
}
