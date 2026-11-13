package setup

import "encoding/json"

type SettingInput struct {
	Value json.RawMessage `json:"value"`
}
