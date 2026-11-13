package team

import "net/http"

type Routes struct {
	Teams Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Teams
	mux.Handle("GET /v1/team/teams", MineHandler{service: s})
	mux.Handle("GET /v1/team/invitations", InvitationsHandler{service: s})
	mux.Handle("GET /v1/team/scopes/{kind}/{id}", ScopeHandler{service: s})
	mux.Handle("POST /v1/team/scopes/{kind}/{id}/teams", CreateHandler{service: s})
	mux.Handle("PUT /v1/team/scopes/{kind}/{id}/seeking", SeekHandler{service: s})
	mux.Handle("POST /v1/team/scopes/{kind}/{id}/distribution", DistributeHandler{service: s})
	mux.Handle("GET /v1/team/teams/{id}", TeamHandler{service: s})
	mux.Handle("PATCH /v1/team/teams/{id}", RenameHandler{service: s})
	mux.Handle("PUT /v1/team/teams/{id}/label", LabelHandler{service: s})
	mux.Handle("POST /v1/team/teams/{id}/invitations", InviteHandler{service: s})
	mux.Handle("DELETE /v1/team/teams/{id}/invitations/{userId}", UninviteHandler{service: s})
	mux.Handle("POST /v1/team/teams/{id}/accept", AcceptHandler{service: s})
	mux.Handle("POST /v1/team/teams/{id}/leave", LeaveHandler{service: s})
	mux.Handle("PUT /v1/team/teams/{id}/members/{userId}", AddMemberHandler{service: s})
	mux.Handle("DELETE /v1/team/teams/{id}/members/{userId}", RemoveMemberHandler{service: s})
	mux.Handle("GET /v1/team/codes/{code}", PreviewHandler{service: s})
	mux.Handle("POST /v1/team/codes/{code}", JoinHandler{service: s})
	mux.Handle("POST /v1/team/leave-requests/{id}", DecideHandler{service: s})
	mux.Handle("POST /v1/team/teams/{id}/proposals", ProposeHandler{service: s})
	mux.Handle("POST /v1/team/proposals/{id}/votes", VoteHandler{service: s})
	mux.Handle("GET /v1/team/card", CardHandler{service: s})
	mux.Handle("PUT /v1/team/card", SaveCardHandler{service: s})
	mux.Handle("POST /v1/team/scopes/{kind}/{id}/suggestions", SuggestHandler{service: s})
}
