package schema

import "encoding/json"

type Value struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
	Set   bool            `json:"set"`
}

func (c Catalog) Resolve(stored map[string]json.RawMessage) []Value {
	values := make([]Value, 0, len(c.list))
	for _, d := range c.list {
		v := Value{Key: d.Key, Value: d.DefaultValue()}
		if raw, ok := stored[d.Key]; ok && d.Validate(raw) == nil {
			v.Value, v.Set = raw, true
		}
		values = append(values, v)
	}
	return values
}
