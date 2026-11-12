package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	DefaultRoleName = "Все участники"
	AdminRoleName   = "Администратор"
)

type Bootstrap struct {
	tenants     database.TenantTx
	roles       RoleStore
	assignments AssignmentStore
}

func NewBootstrap(tenants database.TenantTx) Bootstrap {
	return Bootstrap{tenants: tenants}
}

func (b Bootstrap) EnsureDefaults(ctx context.Context, tenantID string) error {
	return b.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		return b.ensure(ctx, tx)
	})
}

func (b Bootstrap) GrantTopRole(ctx context.Context, tenantID, userID string) (Role, error) {
	var top Role
	err := b.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		if err := b.ensure(ctx, tx); err != nil {
			return err
		}
		roles, err := b.roles.List(ctx, tx)
		if err != nil {
			return err
		}
		top = roles[0]
		if top.Default {
			return errors.New("tenant has only the default role")
		}
		_, err = b.assignments.Upsert(ctx, tx, AssignInput{RoleID: top.ID, UserID: userID}, TenantScope(tenantID), nil)
		return err
	})
	return top, err
}

func (b Bootstrap) ensure(ctx context.Context, tx pgx.Tx) error {
	if _, err := b.roles.Lock(ctx, tx, ""); err != nil {
		return err
	}
	n, err := b.roles.Count(ctx, tx)
	if err != nil || n > 0 {
		return err
	}
	if _, err := b.roles.Insert(ctx, tx, RoleInput{Name: DefaultRoleName}, 0); err != nil {
		return fmt.Errorf("create default role: %w", err)
	}
	if _, err := b.roles.Insert(ctx, tx, RoleInput{Name: AdminRoleName, Color: "#C0392B", Permissions: []Permission{Administrator}}, 1); err != nil {
		return fmt.Errorf("create admin role: %w", err)
	}
	return nil
}
