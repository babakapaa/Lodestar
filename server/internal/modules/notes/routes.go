package notes

import "net/http"

type Routes struct {
	Service Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("GET /v1/notes", ListHandler{service: s})
	mux.Handle("POST /v1/notes", CreateHandler{service: s})
	mux.Handle("GET /v1/notes/{id}", GetHandler{service: s})
	mux.Handle("PUT /v1/notes/{id}/meta", MetaHandler{service: s})
	mux.Handle("DELETE /v1/notes/{id}", DeleteHandler{service: s, deleted: true})
	mux.Handle("POST /v1/notes/{id}/restore", DeleteHandler{service: s})
	mux.Handle("GET /v1/notes/{id}/updates", ChangesHandler{service: s})
	mux.Handle("POST /v1/notes/{id}/updates", PushHandler{service: s})
	mux.Handle("POST /v1/notes/{id}/versions", CompactHandler{service: s})
	mux.Handle("GET /v1/notes/{id}/versions", VersionsHandler{service: s})
	mux.Handle("GET /v1/notes/{id}/versions/{version}", VersionHandler{service: s})
	mux.Handle("GET /v1/notes/{id}/members", MembersHandler{service: s})
	mux.Handle("PUT /v1/notes/{id}/members/{user}", SetMemberHandler{service: s})
	mux.Handle("DELETE /v1/notes/{id}/members/{user}", RemoveMemberHandler{service: s})
}
