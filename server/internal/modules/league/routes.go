package league

import "net/http"

type Routes struct {
	Leagues Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Leagues
	mux.Handle("GET /v1/league/arenas/{kind}/{id}", ViewHandler{service: s})
	mux.Handle("PUT /v1/league/leagues/{id}/rules", RulesHandler{service: s})
	mux.Handle("PUT /v1/league/leagues/{id}/display", DisplayHandler{service: s})
	mux.Handle("POST /v1/league/leagues/{id}/seasons", SeasonHandler{service: s})
	mux.Handle("POST /v1/league/leagues/{id}/adjustments", AdjustHandler{service: s})
	mux.Handle("GET /v1/league/leagues/{id}/players/{userId}", PlayerHandler{service: s})
	mux.Handle("POST /v1/league/leagues/{id}/disputes", DisputeHandler{service: s})
	mux.Handle("POST /v1/league/transitions/{id}", DecideHandler{service: s})
	mux.Handle("POST /v1/league/disputes/{id}", AnswerHandler{service: s})
}
