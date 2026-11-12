package study

import (
	"encoding/json"
	"fmt"
)

type Clock int

func ParseClock(value string) (Clock, error) {
	var h, m int
	if _, err := fmt.Sscanf(value, "%d:%d", &h, &m); err != nil || h < 0 || h > 23 || m < 0 || m > 59 || len(value) != 5 {
		return 0, ErrInvalidSeries
	}
	return Clock(h*60 + m), nil
}

func (c Clock) String() string {
	return fmt.Sprintf("%02d:%02d", int(c)/60, int(c)%60)
}

func (c Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Clock) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseClock(s)
	if err != nil {
		return err
	}
	*c = parsed
	return nil
}
