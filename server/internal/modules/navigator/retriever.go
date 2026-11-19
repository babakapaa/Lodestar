package navigator

import (
	"context"
	"log/slog"
	"sort"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

const (
	searchPool    = 20
	fragmentLimit = 5
)

type Retriever struct {
	tenants  database.TenantTx
	embedder llm.Embedder
	reranker llm.Reranker
	passages PassageStore
}

func (r Retriever) Find(ctx context.Context, tenantID, courseID, query string, match int) []Fragment {
	if r.embedder == nil || courseID == "" || query == "" {
		return nil
	}
	vectors, err := r.embedder.Embed(ctx, []string{query})
	if err != nil {
		slog.Warn("embed navigator question", "error", err)
		return nil
	}
	var found []Fragment
	err = r.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		found, err = r.passages.Search(ctx, tx, r.embedder.Name(), courseID, vectors[0], searchPool, distance(match))
		return err
	})
	if err != nil {
		slog.Warn("search navigator materials", "error", err)
		return nil
	}
	return r.rerank(ctx, query, found)
}

func (r Retriever) rerank(ctx context.Context, query string, found []Fragment) []Fragment {
	if r.reranker != nil && len(found) > 1 {
		texts := make([]string, len(found))
		for i, f := range found {
			texts[i] = f.Text
		}
		scores, err := r.reranker.Rerank(ctx, query, texts)
		if err == nil {
			order := make([]int, len(found))
			for i := range order {
				order[i] = i
			}
			sort.SliceStable(order, func(a, b int) bool { return scores[order[a]] > scores[order[b]] })
			ranked := make([]Fragment, len(found))
			for i, j := range order {
				ranked[i] = found[j]
			}
			found = ranked
		} else {
			slog.Warn("rerank navigator fragments", "error", err)
		}
	}
	if len(found) > fragmentLimit {
		found = found[:fragmentLimit]
	}
	return found
}

func distance(match int) float64 {
	if match <= 0 {
		return 2
	}
	return 1 - float64(match)/100
}
