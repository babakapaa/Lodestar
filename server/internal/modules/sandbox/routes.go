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
	mux.Handle("GET /v1/sandbox/templates", TemplatesHandler{service: s})
	mux.Handle("POST /v1/sandbox/templates", CreateTemplateHandler{service: s})
	mux.Handle("PUT /v1/sandbox/templates/{id}", UpdateTemplateHandler{service: s})
	mux.Handle("GET /v1/sandbox/nodes", NodesHandler{service: s})
	mux.Handle("POST /v1/sandbox/nodes", CreateNodeHandler{service: s})
	mux.Handle("DELETE /v1/sandbox/nodes/{id}", DeleteNodeHandler{service: s})
	mux.Handle("GET /v1/sandbox/runs", RunsHandler{service: s})
	mux.Handle("GET /v1/sandbox/runs/held", HeldRunsHandler{service: s})
	mux.Handle("GET /v1/sandbox/runs/export", ExportRunsHandler{service: s})
	mux.Handle("POST /v1/sandbox/estimate", EstimateHandler{service: s})
	mux.Handle("POST /v1/sandbox/runs/{id}/rerun", RerunHandler{service: s})
	mux.Handle("POST /v1/sandbox/runs/{id}/hold-decision", DecideHoldHandler{service: s})
	mux.Handle("POST /v1/sandbox/runs", StartRunHandler{service: s})
	mux.Handle("GET /v1/sandbox/runs/{id}", RunHandler{service: s})
	mux.Handle("POST /v1/sandbox/runs/{id}/stop", StopRunHandler{service: s})
}

func (r Routes) MountOpen(mux *http.ServeMux) {
	mux.Handle("POST /v1/open/sandbox/agent/sync", AgentHandler{service: r.Quotas})
}
