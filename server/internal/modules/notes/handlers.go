package notes

import (
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ListHandler struct{ service Service }

func (h ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type CreateHandler struct{ service Service }

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	n, created, err := h.service.Create(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, n)
}

type GetHandler struct{ service Service }

func (h GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, n)
}

type MetaHandler struct{ service Service }

func (h MetaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in Meta
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	n, err := h.service.SetMeta(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, n)
}

type DeleteHandler struct {
	service Service
	deleted bool
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.SetDeleted(r.Context(), principal(r), r.PathValue("id"), h.deleted); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type ChangesHandler struct{ service Service }

func (h ChangesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var after int64
	if raw := r.URL.Query().Get("after"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v < 0 {
			badRequest(w)
			return
		}
		after = v
	}
	c, err := h.service.Changes(r.Context(), principal(r), r.PathValue("id"), after)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}

type PushHandler struct{ service Service }

func (h PushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in UpdateInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	u, created, err := h.service.Push(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, map[string]int64{"seq": u.Seq})
}

type CompactHandler struct{ service Service }

func (h CompactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SnapshotInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.Compact(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, v)
}

type VersionsHandler struct{ service Service }

func (h VersionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Versions(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type VersionHandler struct{ service Service }

func (h VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.Version(r.Context(), principal(r), r.PathValue("id"), r.PathValue("version"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, v)
}

type MembersHandler struct{ service Service }

func (h MembersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Members(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type SetMemberHandler struct{ service Service }

func (h SetMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MemberInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.SetMember(r.Context(), principal(r), r.PathValue("id"), r.PathValue("user"), in); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type RemoveMemberHandler struct{ service Service }

func (h RemoveMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.RemoveMember(r.Context(), principal(r), r.PathValue("id"), r.PathValue("user")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
