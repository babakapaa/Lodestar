package apierror

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeNotFound = "not_found"
	CodeInternal = "internal"
)

type Error struct {
	Code string `json:"code"`
}

func Write(w http.ResponseWriter, status int, code string) {
	jsonhttp.Write(w, status, Error{Code: code})
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	Write(w, http.StatusNotFound, CodeNotFound)
}
