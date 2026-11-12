package jsonhttp

import (
	"encoding/json"
	"net/http"
)

const maxBody = 64 << 10

const LargeBody = 1 << 20

func Read(w http.ResponseWriter, r *http.Request, v any) error {
	return ReadLimit(w, r, v, maxBody)
}

func ReadLimit(w http.ResponseWriter, r *http.Request, v any, limit int64) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, limit))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
