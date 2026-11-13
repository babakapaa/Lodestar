package navigator

import "net/http"

type Routes struct {
	Service Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("GET /v1/navigator/dialogs", DialogsHandler{service: s})
	mux.Handle("POST /v1/navigator/dialogs", CreateHandler{service: s})
	mux.Handle("GET /v1/navigator/dialogs/{id}", DialogHandler{service: s})
	mux.Handle("DELETE /v1/navigator/dialogs/{id}", DeleteHandler{service: s})
	mux.Handle("POST /v1/navigator/dialogs/{id}/messages", AskHandler{service: s})
}
