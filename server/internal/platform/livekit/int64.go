package livekit

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Int64 int64

func (n *Int64) UnmarshalJSON(b []byte) error {
	var v int64
	if err := json.Unmarshal([]byte(strings.Trim(string(b), `"`)), &v); err != nil {
		return err
	}
	*n = Int64(v)
	return nil
}

func (n Int64) String() string {
	return strconv.FormatInt(int64(n), 10)
}
