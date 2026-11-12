package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	checkViolation      = "23514"
	foreignKeyViolation = "23503"
	invalidText         = "22P02"
)

const roleColumns = `id::text, name, color, icon, position, permissions, is_default`

type RoleStore struct{}

func (RoleStore) List(ctx context.Context, tx pgx.Tx) ([]Role, error) {
	rows, err := tx.Query(ctx, `SELECT `+roleColumns+` FROM core.roles ORDER BY position DESC`)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	return pgx.CollectRows(rows, scanRole)
}

func (RoleStore) ByIDs(ctx context.Context, tx pgx.Tx, ids []string) ([]Role, error) {
	rows, err := tx.Query(ctx, `SELECT `+roleColumns+` FROM core.roles WHERE id::text = ANY($1) OR is_default`, ids)
	if err != nil {
		return nil, fmt.Errorf("load roles: %w", err)
	}
	return pgx.CollectRows(rows, scanRole)
}

func (RoleStore) Lock(ctx context.Context, tx pgx.Tx, id string) (Role, error) {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended('core.roles:' || current_setting('lodestar.tenant'), 0))`); err != nil {
		return Role{}, fmt.Errorf("lock roles: %w", err)
	}
	if id == "" {
		return Role{}, nil
	}
	rows, err := tx.Query(ctx, `SELECT `+roleColumns+` FROM core.roles WHERE id::text = $1`, id)
	if err != nil {
		return Role{}, fmt.Errorf("find role: %w", err)
	}
	role, err := pgx.CollectExactlyOneRow(rows, scanRole)
	if errors.Is(err, pgx.ErrNoRows) {
		return Role{}, ErrRoleNotFound
	}
	return role, err
}

func (RoleStore) Count(ctx context.Context, tx pgx.Tx) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `SELECT count(*) FROM core.roles`).Scan(&n)
	return n, err
}

func (RoleStore) Insert(ctx context.Context, tx pgx.Tx, in RoleInput, position int) (Role, error) {
	if _, err := tx.Exec(ctx, `UPDATE core.roles SET position = position + 1 WHERE position >= $1`, position); err != nil {
		return Role{}, fmt.Errorf("shift roles: %w", err)
	}
	rows, err := tx.Query(ctx, `
		INSERT INTO core.roles (name, color, icon, position, permissions, is_default)
		VALUES ($1, $2, $3, $4, $5, $4 = 0)
		RETURNING `+roleColumns,
		in.Name, in.Color, in.Icon, position, permissionsOrEmpty(in.Permissions))
	if err != nil {
		return Role{}, roleError(err)
	}
	role, err := pgx.CollectExactlyOneRow(rows, scanRole)
	return role, roleError(err)
}

func (RoleStore) Update(ctx context.Context, tx pgx.Tx, id string, in RoleInput) (Role, error) {
	rows, err := tx.Query(ctx, `
		UPDATE core.roles SET name = $2, color = $3, icon = $4, permissions = $5, updated_at = now()
		WHERE id::text = $1
		RETURNING `+roleColumns,
		id, in.Name, in.Color, in.Icon, permissionsOrEmpty(in.Permissions))
	if err != nil {
		return Role{}, roleError(err)
	}
	role, err := pgx.CollectExactlyOneRow(rows, scanRole)
	return role, roleError(err)
}

func (RoleStore) Delete(ctx context.Context, tx pgx.Tx, r Role) error {
	if _, err := tx.Exec(ctx, `DELETE FROM core.roles WHERE id::text = $1`, r.ID); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	_, err := tx.Exec(ctx, `UPDATE core.roles SET position = position - 1 WHERE position > $1`, r.Position)
	return err
}

func (RoleStore) Move(ctx context.Context, tx pgx.Tx, r Role, to int) error {
	var err error
	switch {
	case to > r.Position:
		_, err = tx.Exec(ctx, `UPDATE core.roles SET position = position - 1 WHERE position > $1 AND position <= $2`, r.Position, to)
	case to < r.Position:
		_, err = tx.Exec(ctx, `UPDATE core.roles SET position = position + 1 WHERE position >= $1 AND position < $2`, to, r.Position)
	default:
		return nil
	}
	if err != nil {
		return fmt.Errorf("shift roles: %w", err)
	}
	_, err = tx.Exec(ctx, `UPDATE core.roles SET position = $2, updated_at = now() WHERE id::text = $1`, r.ID, to)
	return err
}

func scanRole(row pgx.CollectableRow) (Role, error) {
	var r Role
	err := row.Scan(&r.ID, &r.Name, &r.Color, &r.Icon, &r.Position, &r.Permissions, &r.Default)
	return r, err
}

func roleError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == checkViolation {
		return ErrInvalidRole
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrRoleNotFound
	}
	if err != nil {
		return fmt.Errorf("save role: %w", err)
	}
	return nil
}

func permissionsOrEmpty(perms []Permission) []Permission {
	if perms == nil {
		return []Permission{}
	}
	return perms
}
