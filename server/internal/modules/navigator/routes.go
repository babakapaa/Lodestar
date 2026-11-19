package navigator

import "net/http"

type Routes struct {
	Service   Service
	Materials MaterialService
	Topics    TopicService
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("GET /v1/navigator/dialogs", DialogsHandler{service: s})
	mux.Handle("POST /v1/navigator/dialogs", CreateHandler{service: s})
	mux.Handle("GET /v1/navigator/dialogs/{id}", DialogHandler{service: s})
	mux.Handle("DELETE /v1/navigator/dialogs/{id}", DeleteHandler{service: s})
	mux.Handle("POST /v1/navigator/dialogs/{id}/messages", AskHandler{service: s})
	mux.Handle("GET /v1/navigator/journal", JournalHandler{service: s})
	mux.Handle("DELETE /v1/navigator/journal/{id}", ForgetHandler{service: s})
	mux.Handle("GET /v1/navigator/support", SupportHandler{service: s})
	mux.Handle("GET /v1/navigator/help", HelpHandler{service: s})
	mux.Handle("GET /v1/navigator/courses/{id}/materials", MaterialsHandler{service: r.Materials})
	mux.Handle("POST /v1/navigator/courses/{id}/materials", AddMaterialHandler{service: r.Materials})
	mux.Handle("DELETE /v1/navigator/materials/{id}", DeleteMaterialHandler{service: r.Materials})
	mux.Handle("GET /v1/navigator/courses/{id}/topics", TopicsHandler{service: r.Topics})
}
