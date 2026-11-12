package jsonhttp

import (
	"encoding/json"
	"net/http"
)

const maxBody = 64 << 10

func Read(w http.ResponseWriter, r *http.Request, v any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
