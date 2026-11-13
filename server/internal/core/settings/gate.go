package settings

import (
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

const CodeNotLaunched = "university_not_launched"

type Gate struct {
	service       Service
	authenticator auth.Authenticator
	launch        LaunchCheck
}

func NewGate(service Service, authenticator auth.Authenticator, launch LaunchCheck) Gate {
	return Gate{service: service, authenticator: authenticator, launch: launch}
}

func (g Gate) Wrap(name string, h http.Handler) http.Handler {
	return g.authenticator.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, _ := auth.PrincipalFrom(r.Context())
		enabled, err := g.service.ModuleEnabled(r.Context(), p.TenantID, name)
		if err != nil {
			slog.Error("check module state", "module", name, "error", err)
			apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
			return
		}
		if !enabled {
			apierror.NotFound(w, r)
			return
		}
		allowed, err := g.launch.Allowed(r.Context(), p)
		if err != nil {
			slog.Error("check launch", "error", err)
			apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
			return
		}
		if !allowed {
			apierror.Write(w, http.StatusForbidden, CodeNotLaunched)
			return
		}
		h.ServeHTTP(w, r)
	}))
}
