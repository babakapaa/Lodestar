package llm

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const maxReply = 4 << 20

type Client struct {
	transport transport
}

func NewClient(cfg Config) Client {
	return Client{transport: newTransport(cfg)}
}

func (c Client) Complete(ctx context.Context, req Request) (string, error) {
	var reply chatReply
	if err := c.transport.call(ctx, "/chat/completions", c.body(req, false), &reply); err != nil {
		return "", err
	}
	return reply.text(false), nil
}

func (c Client) Stream(ctx context.Context, req Request, delta func(string) error) (string, error) {
	resp, err := c.transport.open(ctx, "/chat/completions", c.body(req, true))
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

func (c Client) body(req Request, stream bool) chatBody {
	body := chatBody{
		Model:       c.transport.cfg.Model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      stream,
	}
	if req.JSON {
		body.ResponseFormat = &responseFormat{Type: "json_object"}
	}
	return body
}
