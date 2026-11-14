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
	mux.Handle("GET /v1/setup/hardware", require(HardwareHandler{check: s.Hardware}))
	mux.Handle("POST /v1/setup/hardware/check", require(HardwareHandler{check: s.CheckHardware}))
	mux.Handle("GET /v1/setup/pilot/courses", require(PilotCoursesHandler{service: s}))
	mux.Handle("POST /v1/setup/launch", require(LaunchHandler{service: s}))
	mux.Handle("POST /v1/setup/expand", require(StepHandler{step: s.Expand}))
}
