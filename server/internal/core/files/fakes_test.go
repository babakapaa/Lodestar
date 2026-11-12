package files_test

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type memoryStorage struct {
	mu      sync.Mutex
	objects map[string]int64
	signed  []string
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{objects: map[string]int64{}}
}

func (m *memoryStorage) PresignPut(_ context.Context, key, contentType string, _ int64, ttl time.Duration) (objectstore.Upload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.signed = append(m.signed, key)
	return objectstore.Upload{URL: "https://s3/" + key, Method: "PUT", Headers: map[string]string{"Content-Type": contentType}, ExpiresAt: time.Now().Add(ttl)}, nil
}

func (m *memoryStorage) PresignGet(_ context.Context, key, _, _ string, _ time.Duration) (string, error) {
	return "https://s3/" + key + "?get", nil
}

func (m *memoryStorage) Size(_ context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	size, ok := m.objects[key]
	if !ok {
		return 0, objectstore.ErrMissing
	}
	return size, nil
}

func (m *memoryStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	size, ok := m.objects[key]
	if !ok {
		return nil, objectstore.ErrMissing
	}
	return io.NopCloser(strings.NewReader(strings.Repeat("x", int(size)))), nil
}

func (m *memoryStorage) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.objects, key)
	return nil
}

func (m *memoryStorage) put(key string, size int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[key] = size
}

func (m *memoryStorage) has(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.objects[key]
	return ok
}

type fixedLimit struct {
	mb int
}

func (f fixedLimit) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	return json.RawMessage(strconv.Itoa(f.mb)), nil
}
