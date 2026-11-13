package llm

import (
	"context"
	"fmt"
)

type EmbedClient struct {
	transport transport
}

func NewEmbedClient(cfg Config) EmbedClient {
	return EmbedClient{transport: newTransport(cfg)}
}

func (c EmbedClient) Name() string {
	return c.transport.cfg.Model
}

func (c EmbedClient) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	var reply embedReply
	if err := c.transport.call(ctx, "/embeddings", embedBody{Model: c.transport.cfg.Model, Input: texts}, &reply); err != nil {
		return nil, err
	}
	out := make([][]float32, len(texts))
	for _, d := range reply.Data {
		if d.Index < 0 || d.Index >= len(out) || len(d.Embedding) == 0 {
			return nil, fmt.Errorf("%w: bad embedding %d", ErrModelFailed, d.Index)
		}
		out[d.Index] = d.Embedding
	}
	for i, v := range out {
		if v == nil {
			return nil, fmt.Errorf("%w: missing embedding %d", ErrModelFailed, i)
		}
	}
	return out, nil
}
