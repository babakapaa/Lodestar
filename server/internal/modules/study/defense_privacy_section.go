package study

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type DefensePrivacySection struct{}

func (DefensePrivacySection) Name() string {
	return "study.defenses"
}

func (DefensePrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT a.title, d.status, s.starts_at, d.held_at,
		       CASE WHEN d.status = 'planned' THEN '[]'::jsonb ELSE d.questions END,
		       CASE WHEN d.status = 'planned' THEN '' ELSE d.summary END
		FROM study.defenses d
		JOIN study.assignments a ON a.id = d.assignment_id
		LEFT JOIN study.defense_slots s ON s.id = d.slot_id
		WHERE d.student_id::text = $1 ORDER BY d.updated_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export defenses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var title, status, summary string
		var startsAt, heldAt *time.Time
		var questions []DefenseQuestion
		err := row.Scan(&title, &status, &startsAt, &heldAt, &questions, &summary)
		return map[string]any{
			"assignment": title, "status": status, "startsAt": startsAt, "heldAt": heldAt,
			"questions": questions, "summary": summary,
		}, err
	})
}
