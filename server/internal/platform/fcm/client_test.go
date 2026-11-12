package fcm_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
)

func TestPushSendsNotification(t *testing.T) {
	var got map[string]map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/projects/campus/messages:send" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"name":"projects/campus/messages/1"}`))
	}))
	defer srv.Close()

	c := fcm.NewWithHTTP(srv.Client(), srv.URL, "campus")
	err := c.Push(context.Background(), "tok", fcm.Message{Title: "Срок", Body: "Завтра", Data: map[string]string{"link": "/study"}})
	if err != nil {
		t.Fatalf("push: %v", err)
	}
	if got["message"]["token"] != "tok" || got["message"]["notification"].(map[string]any)["title"] != "Срок" {
		t.Fatalf("payload = %+v", got)
	}
}

func TestPushReportsUnregisteredToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"status":"INVALID_ARGUMENT","details":[{"errorCode":"UNREGISTERED"}]}}`))
	}))
	defer srv.Close()

	err := fcm.NewWithHTTP(srv.Client(), srv.URL, "campus").Push(context.Background(), "old", fcm.Message{})
	if !errors.Is(err, fcm.ErrUnregistered) {
		t.Fatalf("err = %v, want ErrUnregistered", err)
	}
}

func TestPushReportsOtherFailures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	err := fcm.NewWithHTTP(srv.Client(), srv.URL, "campus").Push(context.Background(), "tok", fcm.Message{})
	if err == nil || errors.Is(err, fcm.ErrUnregistered) {
		t.Fatalf("err = %v", err)
	}
}

func TestNewRejectsIncompleteServiceAccount(t *testing.T) {
	if _, err := fcm.New([]byte(`{"project_id":"campus"}`)); err == nil {
		t.Fatal("incomplete service account accepted")
	}
}
