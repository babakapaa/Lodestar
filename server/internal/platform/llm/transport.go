package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type transport struct {
	cfg  Config
	http *http.Client
}

func newTransport(cfg Config) transport {
	return transport{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

func (t transport) open(ctx context.Context, path string, body any) (*http.Response, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return t.post(ctx, path, "application/json", bytes.NewReader(raw))
}

func (t transport) send(ctx context.Context, path, contentType string, body io.Reader, out any) error {
	resp, err := t.post(ctx, path, contentType, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return decode(resp.Body, out)
}

func (t transport) post(ctx context.Context, path, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(t.cfg.URL, "/")+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if t.cfg.Key != "" {
		req.Header.Set("Authorization", "Bearer "+t.cfg.Key)
	}
	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrModelFailed, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%w: status %d", ErrModelFailed, resp.StatusCode)
	}
	return resp, nil
}

func (t transport) call(ctx context.Context, path string, body, out any) error {
	resp, err := t.open(ctx, path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return decode(resp.Body, out)
}

func decode(body io.Reader, out any) error {
	if err := json.NewDecoder(io.LimitReader(body, maxReply)).Decode(out); err != nil {
		return fmt.Errorf("%w: decode reply: %w", ErrModelFailed, err)
	}
	return nil
}
