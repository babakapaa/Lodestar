package profile

import "encoding/json"

type ValueInput struct {
	Value json.RawMessage `json:"value"`
}
