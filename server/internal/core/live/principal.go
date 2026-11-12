package live

import (
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func principal(r *http.Request) auth.Principal {
	p, _ := auth.PrincipalFrom(r.Context())
	return p
}

func afterParam(r *http.Request) (int64, bool) {
	value := r.URL.Query().Get("after")
	if value == "" {
		return 0, true
	}
	after, err := strconv.ParseInt(value, 10, 64)
	return after, err == nil && after >= 0
}
