package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

type CoursesSource struct{}

func (CoursesSource) Kind() string {
	return "study.courses"
}

func (CoursesSource) Search(ctx context.Context, tx pgx.Tx, p auth.Principal, q search.Query, limit int) ([]search.Hit, error) {
	units, err := structure.Directory{}.UnitsAbove(ctx, tx, p.UserID)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, `
		SELECT c.id::text, c.title, coalesce(c.code, ''),
		       greatest(word_similarity($1, core.search_norm(c.title)),
		                word_similarity($1, core.search_norm(coalesce(c.code, '')))) AS score
		FROM study.courses c
		WHERE c.archived_at IS NULL
		  AND ($1 <% core.search_norm(c.title) OR $1 <% core.search_norm(coalesce(c.code, '')))
		  AND (EXISTS (SELECT 1 FROM study.course_teachers t WHERE t.course_id = c.id AND t.user_id::text = $2)
		       OR EXISTS (SELECT 1 FROM study.course_students s WHERE s.course_id = c.id AND s.user_id::text = $2)
		       OR EXISTS (SELECT 1 FROM study.course_units u WHERE u.course_id = c.id AND u.unit_id::text = ANY($3)))
		ORDER BY score DESC, c.title LIMIT $4`, q.Text, p.UserID, units, limit)
	if err != nil {
		return nil, fmt.Errorf("search courses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (search.Hit, error) {
		var h search.Hit
		err := row.Scan(&h.ID, &h.Title, &h.Subtitle, &h.Score)
		h.Link = "/study/courses/" + h.ID
		return h, err
	})
}
