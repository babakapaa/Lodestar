package delivery

import "encoding/json"

type Item struct {
	Type   string
	Params json.RawMessage
	Link   string
}
