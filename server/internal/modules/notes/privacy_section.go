package notes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	sealer Sealer
}

func (PrivacySection) Name() string {
	return "notes.notes"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT n.meta, n.created_at, n.changed_at, n.deleted_at, v.body
		FROM notes.notes n
		LEFT JOIN LATERAL (SELECT body FROM notes.versions WHERE note_id = n.id ORDER BY seq DESC LIMIT 1) v ON true
		WHERE n.owner_id::text = $1 ORDER BY n.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export notes: %w", err)
	}
	type sealedNote struct {
		meta, body       []byte
		created, changed any
		deleted          any
	}
	sealed, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (sealedNote, error) {
		var n sealedNote
		err := row.Scan(&n.meta, &n.created, &n.changed, &n.deleted, &n.body)
		return n, err
	})
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(sealed))
	for _, n := range sealed {
		meta, err := s.sealer.openMeta(n.meta)
		if err != nil {
			return nil, err
		}
		body := ""
		if n.body != nil {
			plain, err := s.sealer.open(n.body)
			if err != nil {
				return nil, err
			}
			body = string(plain)
		}
		out = append(out, map[string]any{
			"path": meta.Path, "tags": meta.Tags, "links": meta.Links, "body": body,
			"createdAt": n.created, "changedAt": n.changed, "deletedAt": n.deleted,
		})
	}
	return out, nil
}
