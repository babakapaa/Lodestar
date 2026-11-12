package audit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Log struct{}

func (Log) Record(ctx context.Context, tx pgx.Tx, actorID string, c Change) error {
	details, err := json.Marshal(c.Details)
	if err != nil {
		return fmt.Errorf("encode audit details: %w", err)
	}
	if c.Details == nil {
		details = []byte(`{}`)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO core.audit_log (actor_id, action, target_kind, target_id, details)
		VALUES (nullif($1, '')::uuid, $2, $3, $4, $5)`, actorID, c.Action, c.TargetKind, c.TargetID, details)
	if err != nil {
		return fmt.Errorf("record %s: %w", c.Action, err)
	}
	return nil
}
