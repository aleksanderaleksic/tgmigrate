package common

import (
	"fmt"
	"time"
)

type JSONTime time.Time

const dateTimeFormat = `"` + time.RFC3339 + `"`

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(dateTimeFormat))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(value []byte) error {
	ti, err := time.Parse(dateTimeFormat, string(value))
	if err != nil {
		return err
	}
	*t = JSONTime(ti)
	return nil
}

func (t JSONTime) Time() time.Time {
	return time.Time(t)
}
