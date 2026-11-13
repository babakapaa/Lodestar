package chat

import "net/http"

type Routes struct {
	Chats      Service
	Messages   MessageService
	Blocks     BlockService
	Reports    ReportService
	Calls      CallService
	Recordings RecordingService
}

func (r Routes) Mount(mux *http.ServeMux) {
	s, m := r.Chats, r.Messages
	mux.Handle("GET /v1/chat/chats", ChatsHandler{service: s})
	mux.Handle("POST /v1/chat/chats", CreateHandler{service: s})
	mux.Handle("POST /v1/chat/direct", DirectHandler{service: s})
	mux.Handle("GET /v1/chat/chats/{id}", ChatHandler{service: s})
	mux.Handle("PATCH /v1/chat/chats/{id}", UpdateHandler{service: s})
	mux.Handle("POST /v1/chat/chats/{id}/members", AddMembersHandler{service: s})
	mux.Handle("DELETE /v1/chat/chats/{id}/members/{userId}", RemoveMemberHandler{service: s})
	mux.Handle("PUT /v1/chat/chats/{id}/mute", MuteHandler{service: s})
	mux.Handle("POST /v1/chat/chats/{id}/accept", AcceptHandler{service: s})
	mux.Handle("POST /v1/chat/chats/{id}/read", ReadHandler{service: s})
	mux.Handle("GET /v1/chat/chats/{id}/messages", MessagesHandler{service: m})
	mux.Handle("POST /v1/chat/chats/{id}/messages", SendHandler{service: m})
	mux.Handle("PATCH /v1/chat/messages/{id}", EditHandler{service: m})
	mux.Handle("DELETE /v1/chat/messages/{id}", DeleteMessageHandler{service: m})
	mux.Handle("GET /v1/chat/blocks", BlocksHandler{service: r.Blocks})
	mux.Handle("POST /v1/chat/blocks", BlockHandler{service: r.Blocks})
	mux.Handle("DELETE /v1/chat/blocks/{userId}", UnblockHandler{service: r.Blocks})
	mux.Handle("GET /v1/chat/reports", ReportsHandler{service: r.Reports})
	mux.Handle("POST /v1/chat/reports", ReportHandler{service: r.Reports})
	mux.Handle("POST /v1/chat/reports/{id}/resolve", ResolveHandler{service: r.Reports})
	mux.Handle("GET /v1/chat/chats/{id}/call", CallHandler{service: r.Calls})
	mux.Handle("POST /v1/chat/chats/{id}/call", JoinCallHandler{service: r.Calls})
	mux.Handle("POST /v1/chat/calls/{id}/leave", LeaveCallHandler{service: r.Calls})
	mux.Handle("POST /v1/chat/calls/{id}/recording", StartRecordingHandler{service: r.Recordings})
	mux.Handle("DELETE /v1/chat/calls/{id}/recording", StopRecordingHandler{service: r.Recordings})
}

func (r Routes) MountOpen(mux *http.ServeMux) {
	mux.Handle("POST /v1/open/chat/livekit", CallHookHandler{service: r.Calls, recordings: r.Recordings})
}
