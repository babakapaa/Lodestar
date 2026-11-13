package attendance

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/clientip"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CheckHandler struct{ service Service }

func (h CheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in CheckInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	from, _ := clientip.Of(r)
	res, err := h.service.Check(r.Context(), principal(r), in, from)
	respond(w, res, err)
}

type OpenLessonsHandler struct{ service Service }

func (h OpenLessonsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.OpenLessons(r.Context(), principal(r))
	respond(w, list, err)
}

type MineHandler struct{ service Service }

func (h MineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Mine(r.Context(), principal(r))
	respond(w, list, err)
}
