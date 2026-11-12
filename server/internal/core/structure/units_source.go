package structure

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

type UnitsSource struct{}

func (UnitsSource) Kind() string {
	return "units"
}

func (UnitsSource) Search(ctx context.Context, tx pgx.Tx, _ auth.Principal, q search.Query, limit int) ([]search.Hit, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, name, kind, coalesce(code, ''),
		       greatest(word_similarity($1, core.search_norm(name)),
		                word_similarity($1, core.search_norm(coalesce(code, '')))) AS score
		FROM core.units
		WHERE archived_at IS NULL
		  AND ($1 <% core.search_norm(name) OR $1 <% core.search_norm(coalesce(code, '')))
		ORDER BY score DESC, name LIMIT $2`, q.Text, limit)
	if err != nil {
		return nil, fmt.Errorf("search units: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (search.Hit, error) {
		var h search.Hit
		var kind, code string
		err := row.Scan(&h.ID, &h.Title, &kind, &code, &h.Score)
		h.Subtitle = strings.TrimSpace(kind + " " + code)
		h.Link = "/units/" + h.ID
		return h, err
	})
}
