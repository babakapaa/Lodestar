package llm

import (
	"context"
	"fmt"
)

type RerankClient struct {
	transport transport
}

func NewRerankClient(cfg Config) RerankClient {
	return RerankClient{transport: newTransport(cfg)}
}

func (c RerankClient) Rerank(ctx context.Context, query string, documents []string) ([]float64, error) {
	var reply rerankReply
	body := rerankBody{Model: c.transport.cfg.Model, Query: query, Documents: documents}
	if err := c.transport.call(ctx, "/rerank", body, &reply); err != nil {
		return nil, err
	}
	scores := make([]float64, len(documents))
	for _, r := range reply.Results {
		if r.Index < 0 || r.Index >= len(scores) {
			return nil, fmt.Errorf("%w: bad rerank index %d", ErrModelFailed, r.Index)
		}
		scores[r.Index] = r.Score
	}
	return scores, nil
}
