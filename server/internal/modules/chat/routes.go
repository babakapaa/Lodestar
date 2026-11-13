package chat

import "net/http"

type Routes struct {
	Chats    Service
	Messages MessageService
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
	mux.Handle("POST /v1/chat/chats/{id}/read", ReadHandler{service: s})
	mux.Handle("GET /v1/chat/chats/{id}/messages", MessagesHandler{service: m})
	mux.Handle("POST /v1/chat/chats/{id}/messages", SendHandler{service: m})
	mux.Handle("PATCH /v1/chat/messages/{id}", EditHandler{service: m})
	mux.Handle("DELETE /v1/chat/messages/{id}", DeleteMessageHandler{service: m})
}
