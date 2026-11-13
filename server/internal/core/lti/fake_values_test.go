package lti_test

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/jackc/pgx/v5"
)

type fakeValues struct {
	mu     sync.Mutex
	source string
}

func (f *fakeValues) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return json.Marshal(f.source)
}

func (f *fakeValues) set(source string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.source = source
}
