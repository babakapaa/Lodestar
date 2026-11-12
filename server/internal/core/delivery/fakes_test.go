package delivery_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

type values map[string]string

func (v values) get(key string) (json.RawMessage, error) {
	raw, ok := v[key]
	if !ok {
		return nil, errors.New("unexpected key " + key)
	}
	return json.RawMessage(raw), nil
}

type preferences struct {
	values
}

func (p preferences) Value(_ context.Context, _ pgx.Tx, _, key string) (json.RawMessage, error) {
	return p.get(key)
}

type settings struct {
	values
}

func (s settings) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	return s.get(key)
}

type pushed struct {
	token string
	msg   fcm.Message
}

type fakePusher struct {
	mu   sync.Mutex
	sent []pushed
	fail map[string]error
}

func (f *fakePusher) Push(_ context.Context, token string, m fcm.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.fail[token]; err != nil {
		return err
	}
	f.sent = append(f.sent, pushed{token: token, msg: m})
	return nil
}

type fakeMailer struct {
	mu   sync.Mutex
	sent []mail.Message
	err  error
}

func (f *fakeMailer) Send(_ context.Context, m mail.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, m)
	return nil
}
