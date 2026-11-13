package navigator

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TopicStore struct{}

func (TopicStore) Record(ctx context.Context, tx pgx.Tx, courseID, userID string, intent Intent, materials []string, on time.Time) error {
	ids := []*string{nil}
	if len(materials) > 0 {
		ids = make([]*string, len(materials))
		for i := range materials {
			ids[i] = &materials[i]
		}
	}
	for _, id := range ids {
		_, err := tx.Exec(ctx, `
			INSERT INTO navigator.topic_hits (course_id, material_id, user_id, intent, asked_on)
			SELECT $1, m.id, $3, $4, $5::date FROM (SELECT $2::uuid AS id) m
			WHERE m.id IS NULL OR EXISTS (SELECT 1 FROM navigator.materials x WHERE x.id = m.id)`,
			courseID, id, userID, string(intent), on.Format(time.DateOnly))
		if err != nil {
			return fmt.Errorf("record topic: %w", err)
		}
	}
	return nil
}

func (TopicStore) Students(ctx context.Context, tx pgx.Tx, courseID string, from time.Time) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `
		SELECT count(DISTINCT user_id) FROM navigator.topic_hits WHERE course_id = $1 AND asked_on >= $2::date`,
		courseID, from.Format(time.DateOnly)).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count askers: %w", err)
	}
	return n, nil
}

func (TopicStore) Topics(ctx context.Context, tx pgx.Tx, courseID string, from, recent time.Time, least int) ([]TopicCount, error) {
	rows, err := tx.Query(ctx, `
		SELECT h.material_id::text, m.title, count(DISTINCT h.user_id),
		       count(DISTINCT h.user_id) FILTER (WHERE h.asked_on >= $3::date)
		FROM navigator.topic_hits h JOIN navigator.materials m ON m.id = h.material_id
		WHERE h.course_id = $1 AND h.asked_on >= $2::date
		GROUP BY h.material_id, m.title
		HAVING count(DISTINCT h.user_id) >= $4
		ORDER BY 3 DESC, 2 LIMIT 20`,
		courseID, from.Format(time.DateOnly), recent.Format(time.DateOnly), least)
	if err != nil {
		return nil, fmt.Errorf("course topics: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (TopicCount, error) {
		var t TopicCount
		err := row.Scan(&t.MaterialID, &t.Title, &t.Students, &t.Recent)
		return t, err
	})
}

func (TopicStore) Intents(ctx context.Context, tx pgx.Tx, courseID string, from time.Time, least int) ([]IntentCount, error) {
	rows, err := tx.Query(ctx, `
		SELECT intent, count(DISTINCT user_id) FROM navigator.topic_hits
		WHERE course_id = $1 AND asked_on >= $2::date AND intent <> ''
		GROUP BY intent HAVING count(DISTINCT user_id) >= $3 ORDER BY 2 DESC, 1`,
		courseID, from.Format(time.DateOnly), least)
	if err != nil {
		return nil, fmt.Errorf("course intents: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (IntentCount, error) {
		var c IntentCount
		err := row.Scan(&c.Intent, &c.Students)
		return c, err
	})
}

func (TopicStore) Weekly(ctx context.Context, tx pgx.Tx, courseID string, from time.Time) (map[string]int, error) {
	rows, err := tx.Query(ctx, `
		SELECT date_trunc('week', asked_on)::date::text, count(DISTINCT user_id) FROM navigator.topic_hits
		WHERE course_id = $1 AND asked_on >= $2::date GROUP BY 1`,
		courseID, from.Format(time.DateOnly))
	if err != nil {
		return nil, fmt.Errorf("weekly askers: %w", err)
	}
	out := map[string]int{}
	var week string
	var n int
	_, err = pgx.ForEachRow(rows, []any{&week, &n}, func() error {
		out[week] = n
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("weekly askers: %w", err)
	}
	return out, nil
}
