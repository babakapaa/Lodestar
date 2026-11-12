package schema

import (
	"encoding/json"
	"fmt"
	"slices"
	"unicode/utf8"
)

type Definition struct {
	Key       string   `json:"key"`
	Type      Type     `json:"type"`
	Options   []string `json:"options,omitempty"`
	Min       *int     `json:"min,omitempty"`
	Max       *int     `json:"max,omitempty"`
	MaxLength int      `json:"maxLength,omitempty"`
	Default   any      `json:"default"`
}

func (d Definition) Validate(raw json.RawMessage) error {
	var err error
	switch d.Type {
	case Bool:
		var v bool
		err = json.Unmarshal(raw, &v)
	case Integer:
		err = d.validateInteger(raw)
	case Text:
		var v string
		if err = json.Unmarshal(raw, &v); err == nil && d.MaxLength > 0 && utf8.RuneCountInString(v) > d.MaxLength {
			err = fmt.Errorf("longer than %d", d.MaxLength)
		}
	case Choice:
		var v string
		if err = json.Unmarshal(raw, &v); err == nil && !slices.Contains(d.Options, v) {
			err = fmt.Errorf("%q is not an option", v)
		}
	case TimeZone:
		var v string
		if err = json.Unmarshal(raw, &v); err == nil && v != "" {
			err = validZone(v)
		}
	default:
		err = fmt.Errorf("type %q", d.Type)
	}
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrInvalidValue, d.Key, err)
	}
	return nil
}

func (d Definition) validateInteger(raw json.RawMessage) error {
	var v int
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	if d.Min != nil && v < *d.Min {
		return fmt.Errorf("below %d", *d.Min)
	}
	if d.Max != nil && v > *d.Max {
		return fmt.Errorf("above %d", *d.Max)
	}
	return nil
}

func (d Definition) DefaultValue() json.RawMessage {
	raw, err := json.Marshal(d.Default)
	if err != nil {
		return json.RawMessage("null")
	}
	return raw
}

func Bounds(lo, hi int) (*int, *int) {
	return &lo, &hi
}
