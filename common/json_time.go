package common

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(value []byte) error {
	ti, err := time.Parse(`"`+time.RFC3339+`"`, string(value))
	if err != nil {
		return err
	}
	*t = JSONTime(ti)
	return nil
}

func (t JSONTime) Time() time.Time {
	return time.Time(t)
}
