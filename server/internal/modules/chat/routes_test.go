package chat_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestRoutesAnswerWithStatusAndCodes(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("aigerim"), w.user("erlan")

	rec := w.call(a, http.MethodPost, "/v1/chat/direct", `{"userId":"`+b.UserID+`"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("direct: %d %s", rec.Code, rec.Body)
	}
	var d struct{ ID string }
	_ = json.Unmarshal(rec.Body.Bytes(), &d)
	if rec := w.call(a, http.MethodPost, "/v1/chat/direct", `{"userId":"`+b.UserID+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("direct again: %d", rec.Code)
	}
	if rec := w.call(a, http.MethodPost, "/v1/chat/chats/"+d.ID+"/messages", `{"clientId":"c1","text":"Привет"}`); rec.Code != http.StatusCreated {
		t.Fatalf("send: %d %s", rec.Code, rec.Body)
	}
	rec = w.call(b, http.MethodGet, "/v1/chat/chats/"+d.ID+"/messages?since=0", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Привет") {
		t.Fatalf("messages: %d %s", rec.Code, rec.Body)
	}
	if rec := w.call(b, http.MethodGet, "/v1/chat/chats/"+d.ID+"/messages?since=1&before=2", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad query: %d", rec.Code)
	}
	rec = w.call(w.user("dana"), http.MethodGet, "/v1/chat/chats/"+d.ID, "")
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "chat_not_found") {
		t.Fatalf("outsider: %d %s", rec.Code, rec.Body)
	}
	if rec := w.call(a, http.MethodPut, "/v1/chat/chats/"+d.ID+"/mute", `{"muted":true}`); rec.Code != http.StatusOK {
		t.Fatalf("mute: %d", rec.Code)
	}
	if rec := w.call(a, http.MethodPost, "/v1/chat/chats", `{"kind":"direct","title":"x"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("create direct via chats: %d", rec.Code)
	}
}
