package access

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type OverrideStore struct{}

type scopedOverride struct {
	scope    Scope
	override Override
}

func (OverrideStore) InChain(ctx context.Context, tx pgx.Tx, chain []Scope, roleIDs []string) ([]scopedOverride, error) {
	kinds, ids := splitChain(chain)
	rows, err := tx.Query(ctx, `
		SELECT scope_kind, scope_id, role_id::text, allow, deny FROM core.permission_overrides
		WHERE role_id::text = ANY($3)
		  AND (scope_kind, scope_id) IN (SELECT * FROM unnest($1::text[], $2::text[]))`,
		kinds, ids, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("load overrides: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (scopedOverride, error) {
		var so scopedOverride
		err := row.Scan(&so.scope.Kind, &so.scope.ID, &so.override.RoleID, &so.override.Allow, &so.override.Deny)
		return so, err
	})
}

func (OverrideStore) At(ctx context.Context, tx pgx.Tx, s Scope) ([]Override, error) {
	rows, err := tx.Query(ctx, `
		SELECT role_id::text, allow, deny FROM core.permission_overrides
		WHERE scope_kind = $1 AND scope_id = $2`, s.Kind, s.ID)
	if err != nil {
		return nil, fmt.Errorf("list overrides: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Override, error) {
		var o Override
		err := row.Scan(&o.RoleID, &o.Allow, &o.Deny)
		return o, err
	})
}

func (OverrideStore) Put(ctx context.Context, tx pgx.Tx, s Scope, o Override) error {
	if len(o.Allow) == 0 && len(o.Deny) == 0 {
		_, err := tx.Exec(ctx, `
			DELETE FROM core.permission_overrides WHERE scope_kind = $1 AND scope_id = $2 AND role_id::text = $3`,
			s.Kind, s.ID, o.RoleID)
		return err
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO core.permission_overrides (scope_kind, scope_id, role_id, allow, deny)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id, scope_kind, scope_id, role_id) DO UPDATE
		SET allow = EXCLUDED.allow, deny = EXCLUDED.deny`,
		s.Kind, s.ID, o.RoleID, permissionsOrEmpty(o.Allow), permissionsOrEmpty(o.Deny))
	if err != nil {
		return fmt.Errorf("save override: %w", err)
	}
	return nil
}
