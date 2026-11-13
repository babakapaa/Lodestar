package stats

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const CodeInvalidPeriod = "invalid_stats_period"

type UnitsHandler struct {
	service Service
}

func (h UnitsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var parent *string
	if v := q.Get("parent"); v != "" {
		parent = &v
	}
	days := 0
	if v := q.Get("days"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			apierror.Write(w, http.StatusBadRequest, CodeInvalidPeriod)
			return
		}
		days = n
	}
	p, _ := auth.PrincipalFrom(r.Context())
	report, err := h.service.Units(r.Context(), p, parent, days)
	switch {
	case err == nil:
		jsonhttp.Write(w, http.StatusOK, report)
	case errors.Is(err, ErrInvalidPeriod):
		apierror.Write(w, http.StatusBadRequest, CodeInvalidPeriod)
	case errors.Is(err, structure.ErrUnitNotFound):
		apierror.Write(w, http.StatusNotFound, structure.CodeUnitNotFound)
	default:
		access.WriteError(w, err)
	}
}
