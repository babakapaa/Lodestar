package sandbox

import "net/http"

type Routes struct {
	Quotas Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Quotas
	mux.Handle("GET /v1/sandbox/quota", QuotaHandler{service: s})
	mux.Handle("POST /v1/sandbox/transfers", TransferHandler{service: s})
	mux.Handle("GET /v1/sandbox/topups", TopupsHandler{service: s})
	mux.Handle("POST /v1/sandbox/topups", RequestTopupHandler{service: s})
	mux.Handle("POST /v1/sandbox/topups/{id}/decision", DecideTopupHandler{service: s})
	mux.Handle("GET /v1/sandbox/classes", ClassesHandler{service: s})
	mux.Handle("POST /v1/sandbox/classes", CreateClassHandler{service: s})
	mux.Handle("PUT /v1/sandbox/classes/{id}", UpdateClassHandler{service: s})
	mux.Handle("GET /v1/sandbox/terms", TermsHandler{service: s})
	mux.Handle("POST /v1/sandbox/terms", CreateTermHandler{service: s})
	mux.Handle("PUT /v1/sandbox/terms/{id}", UpdateTermHandler{service: s})
	mux.Handle("DELETE /v1/sandbox/terms/{id}", DeleteTermHandler{service: s})
}
