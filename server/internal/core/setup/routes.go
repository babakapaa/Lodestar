package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	require := r.Authenticator.Require
	s := r.Service
	mux.Handle("GET /v1/setup/status", require(StatusHandler{service: s}))
	mux.Handle("GET /v1/setup", require(OverviewHandler{service: s}))
	mux.Handle("PUT /v1/setup/sections/{section}/assignee", require(AssignHandler{service: s}))
	mux.Handle("PUT /v1/setup/sections/{section}/settings/{key}", require(SettingHandler{service: s}))
	mux.Handle("POST /v1/setup/sections/{section}/submit", require(SubmitHandler{service: s}))
	mux.Handle("POST /v1/setup/sections/{section}/review", require(ReviewHandler{service: s}))
	mux.Handle("PUT /v1/setup/legal/{question}", require(AnswerHandler{service: s}))
	mux.Handle("POST /v1/setup/legal/{question}/approve", require(ApproveAnswerHandler{service: s}))
	mux.Handle("GET /v1/setup/protocol", require(ProtocolHandler{service: s}))
	mux.Handle("POST /v1/setup/protocol/sign", require(StepHandler{step: s.Sign}))
	mux.Handle("POST /v1/setup/launch", require(StepHandler{step: s.Launch}))
}
