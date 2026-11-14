package chat_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
)

type fakeSpeech struct {
	mu      sync.Mutex
	objects map[string]string
	text    string
	broken  bool
	heard   []string
}

func newFakeSpeech() *fakeSpeech {
	return &fakeSpeech{objects: map[string]string{}}
}

func (f *fakeSpeech) Transcribe(_ context.Context, audio io.Reader, name string) (string, error) {
	raw, err := io.ReadAll(audio)
	if err != nil {
		return "", err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.heard = append(f.heard, name+":"+string(raw))
	if f.broken {
		return "", errors.New("speech server down")
	}
	return f.text, nil
}

func (f *fakeSpeech) Open(_ context.Context, key string) (io.ReadCloser, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	content, ok := f.objects[key]
	if !ok {
		content = "audio"
	}
	return io.NopCloser(strings.NewReader(content)), nil
}

func (f *fakeSpeech) Put(_ context.Context, key, _ string, body io.Reader, _ int64) error {
	raw, err := io.ReadAll(body)
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objects[key] = string(raw)
	return err
}
