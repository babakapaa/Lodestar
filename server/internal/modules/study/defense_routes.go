package study

import "net/http"

func (r Routes) mountDefenses(mux *http.ServeMux) {
	d := r.Defenses
	mux.Handle("GET /v1/study/assignments/{id}/defense-slots", SlotsHandler{service: d})
	mux.Handle("POST /v1/study/assignments/{id}/defense-slots", AddSlotsHandler{service: d})
	mux.Handle("DELETE /v1/study/defense-slots/{id}", DeleteSlotHandler{service: d})
	mux.Handle("GET /v1/study/assignments/{id}/defense", MyDefenseHandler{service: d})
	mux.Handle("PUT /v1/study/assignments/{id}/defense/booking", BookHandler{service: d})
	mux.Handle("DELETE /v1/study/assignments/{id}/defense/booking", UnbookHandler{service: d})
	mux.Handle("GET /v1/study/assignments/{id}/defenses", DefensesHandler{service: d})
	mux.Handle("POST /v1/study/assignments/{id}/defenses/sample", SampleHandler{service: d})
	mux.Handle("GET /v1/study/assignments/{id}/defenses/{student}", DefenseHandler{service: d})
	mux.Handle("PUT /v1/study/assignments/{id}/defenses/{student}", SaveDefenseHandler{service: d})
	mux.Handle("PUT /v1/study/assignments/{id}/defenses/{student}/selection", SelectHandler{service: d})
	mux.Handle("POST /v1/study/assignments/{id}/defenses/{student}/finish", FinishHandler{service: d})
	mux.Handle("POST /v1/study/assignments/{id}/defenses/{student}/protocol-draft", DraftProtocolHandler{service: d})
}
