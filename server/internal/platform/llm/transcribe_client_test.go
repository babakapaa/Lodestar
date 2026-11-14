package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

func TestTranscribe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/transcriptions" || r.Header.Get("Authorization") != "Bearer k" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		audio, _ := io.ReadAll(file)
		if string(audio) == "broken" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"text": r.FormValue("model") + ":" + header.Filename + ":" + string(audio)})
	}))
	defer srv.Close()
	c := llm.NewTranscribeClient(llm.Config{URL: srv.URL + "/v1", Model: "whisper", Key: "k", Timeout: time.Second})

	got, err := c.Transcribe(context.Background(), strings.NewReader("звук"), "call.ogg")
	if err != nil || got != "whisper:call.ogg:звук" {
		t.Fatalf("transcribe = %q, %v", got, err)
	}
	if _, err := c.Transcribe(context.Background(), strings.NewReader("broken"), "call.ogg"); !errors.Is(err, llm.ErrModelFailed) {
		t.Fatalf("broken = %v", err)
	}
}
