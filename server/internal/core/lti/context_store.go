package lti

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const contextColumns = `id::text, platform_id::text, deployment_id, context_id, title, label, members_url, lineitems_url, ags_scopes, origin_kind, origin_id,
	linked_by::text, assignment_source, synced_at, sync_error`

type ContextStore struct{}

func (ContextStore) Upsert(ctx context.Context, tx pgx.Tx, platformID string, c LaunchClaims) (Context, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO lti.contexts (platform_id, deployment_id, context_id, title, label, members_url, lineitems_url, ags_scopes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (platform_id, deployment_id, context_id) DO UPDATE
		SET title = EXCLUDED.title, label = EXCLUDED.label,
		    members_url = coalesce(EXCLUDED.members_url, lti.contexts.members_url),
		    lineitems_url = coalesce(EXCLUDED.lineitems_url, lti.contexts.lineitems_url),
		    ags_scopes = CASE WHEN cardinality(EXCLUDED.ags_scopes) > 0 THEN EXCLUDED.ags_scopes ELSE lti.contexts.ags_scopes END,
		    updated_at = now()
		RETURNING `+contextColumns,
		platformID, c.DeploymentID, c.Context.ID, c.Context.Title, c.Context.Label,
		optional(c.Names.MembershipsURL), optional(c.Grades.LineItemsURL), nonNil(c.Grades.Scopes))
	if err != nil {
		return Context{}, fmt.Errorf("save lti context: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanContext)
}

func (ContextStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Context, error) {
	rows, err := tx.Query(ctx, `SELECT `+contextColumns+` FROM lti.contexts WHERE id::text = $1`, id)
	if err != nil {
		return Context{}, fmt.Errorf("find lti context: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanContext)
	if errors.Is(err, pgx.ErrNoRows) {
		return Context{}, ErrContextNotFound
	}
	return c, err
}

func (ContextStore) Link(ctx context.Context, tx pgx.Tx, id, kind, originID, userID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE lti.contexts SET origin_kind = $2, origin_id = $3, linked_by = $4, linked_at = now(), updated_at = now()
		WHERE id::text = $1`, id, kind, originID, userID)
	if err != nil {
		return fmt.Errorf("link lti context: %w", err)
	}
	return nil
}

func (ContextStore) Unlink(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `
		UPDATE lti.contexts SET origin_kind = NULL, origin_id = NULL, linked_by = NULL, linked_at = NULL, updated_at = now()
		WHERE id::text = $1`, id)
	if err != nil {
		return fmt.Errorf("unlink lti context: %w", err)
	}
	return nil
}

func (ContextStore) Linked(ctx context.Context, tx pgx.Tx) ([]Context, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+contextColumns+` FROM lti.contexts
		WHERE origin_kind IS NOT NULL AND (members_url IS NOT NULL OR lineitems_url IS NOT NULL)
		ORDER BY synced_at NULLS FIRST, id`)
	if err != nil {
		return nil, fmt.Errorf("linked lti contexts: %w", err)
	}
	return pgx.CollectRows(rows, scanContext)
}

func (ContextStore) ByOrigin(ctx context.Context, tx pgx.Tx, kind, originID string) ([]Context, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+contextColumns+` FROM lti.contexts WHERE origin_kind = $1 AND origin_id = $2 ORDER BY title, id`, kind, originID)
	if err != nil {
		return nil, fmt.Errorf("lti contexts by origin: %w", err)
	}
	return pgx.CollectRows(rows, scanContext)
}

func (ContextStore) SetSource(ctx context.Context, tx pgx.Tx, id string, source *string) error {
	_, err := tx.Exec(ctx, `UPDATE lti.contexts SET assignment_source = $2, updated_at = now() WHERE id::text = $1`, id, source)
	if err != nil {
		return fmt.Errorf("set lti assignment source: %w", err)
	}
	return nil
}

func (ContextStore) Synced(ctx context.Context, tx pgx.Tx, id string, syncErr *string) error {
	_, err := tx.Exec(ctx, `UPDATE lti.contexts SET synced_at = now(), sync_error = $2 WHERE id::text = $1`, id, syncErr)
	if err != nil {
		return fmt.Errorf("mark lti context synced: %w", err)
	}
	return nil
}

func scanContext(row pgx.CollectableRow) (Context, error) {
	var c Context
	err := row.Scan(&c.ID, &c.PlatformID, &c.DeploymentID, &c.ContextID, &c.Title, &c.Label,
		&c.MembersURL, &c.LineItemsURL, &c.Scopes, &c.OriginKind, &c.OriginID,
		&c.LinkedBy, &c.Source, &c.SyncedAt, &c.SyncError)
	return c, err
}

func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nonNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
