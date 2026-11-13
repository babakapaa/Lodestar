package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxReply = 4 << 20

type Client struct {
	cfg  Config
	http *http.Client
}

func NewClient(cfg Config) Client {
	return Client{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

func (c Client) Complete(ctx context.Context, req Request) (string, error) {
	resp, err := c.send(ctx, req, false)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	var reply chatReply
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxReply)).Decode(&reply); err != nil {
		return "", fmt.Errorf("%w: decode reply: %w", ErrModelFailed, err)
	}
	return reply.text(false), nil
}

func (c Client) Stream(ctx context.Context, req Request, delta func(string) error) (string, error) {
	resp, err := c.send(ctx, req, true)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	var full strings.Builder
	scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxReply))
	scanner.Buffer(make([]byte, 64<<10), 1<<20)
	for scanner.Scan() {
		data, ok := strings.CutPrefix(scanner.Text(), "data:")
		data = strings.TrimSpace(data)
		if !ok || data == "" {
			continue
		}
		if data == "[DONE]" {
			break
		}
		var reply chatReply
		if err := json.Unmarshal([]byte(data), &reply); err != nil {
			return "", fmt.Errorf("%w: decode chunk: %w", ErrModelFailed, err)
		}
		piece := reply.text(true)
		if piece == "" {
			continue
		}
		full.WriteString(piece)
		if err := delta(piece); err != nil {
			return "", err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("%w: read stream: %w", ErrModelFailed, err)
	}
	return full.String(), nil
}

func (c Client) send(ctx context.Context, req Request, stream bool) (*http.Response, error) {
	body := chatBody{
		Model:       c.cfg.Model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      stream,
	}
	if req.JSON {
		body.ResponseFormat = &responseFormat{Type: "json_object"}
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.URL, "/")+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.Key != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.cfg.Key)
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrModelFailed, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%w: status %d", ErrModelFailed, resp.StatusCode)
	}
	return resp, nil
}
