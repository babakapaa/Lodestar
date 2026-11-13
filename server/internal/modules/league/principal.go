package league

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func principal(r *http.Request) auth.Principal {
	p, _ := auth.PrincipalFrom(r.Context())
	return p
}
