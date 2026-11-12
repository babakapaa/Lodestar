package profile

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

type PeopleSource struct{}

func (PeopleSource) Kind() string {
	return "people"
}

func (PeopleSource) Search(ctx context.Context, tx pgx.Tx, _ auth.Principal, q search.Query, limit int) ([]search.Hit, error) {
	rows, err := tx.Query(ctx, `
		SELECT u.id::text, coalesce(nullif(p.display_name, ''), u.display_name) AS name,
		       greatest(word_similarity($1, core.search_norm(u.display_name)),
		                word_similarity($1, core.search_norm(coalesce(p.display_name, '')))) AS score
		FROM core.users u LEFT JOIN core.profiles p ON p.user_id = u.id
		WHERE $1 <% core.search_norm(u.display_name) OR $1 <% core.search_norm(coalesce(p.display_name, ''))
		ORDER BY score DESC, name LIMIT $2`, q.Text, limit)
	if err != nil {
		return nil, fmt.Errorf("search people: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (search.Hit, error) {
		var h search.Hit
		err := row.Scan(&h.ID, &h.Title, &h.Score)
		h.Link = "/people/" + h.ID
		return h, err
	})
}
