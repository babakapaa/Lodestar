package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

func readOptional(w http.ResponseWriter, r *http.Request, into any) bool {
	if r.ContentLength == 0 {
		return true
	}
	if err := jsonhttp.Read(w, r, into); err != nil {
		badRequest(w)
		return false
	}
	return true
}
