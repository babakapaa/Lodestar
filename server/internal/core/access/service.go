package access

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants     database.TenantTx
	catalog     Catalog
	resolver    Resolver
	roles       RoleStore
	assignments AssignmentStore
	overrides   OverrideStore
}

func NewService(tenants database.TenantTx, catalog Catalog, tree ScopeTree) Service {
	return Service{tenants: tenants, catalog: catalog, resolver: NewResolver(tree, catalog)}
}

func (s Service) Catalog() []Permission {
	return s.catalog.All()
}

func (s Service) Roles(ctx context.Context, p auth.Principal) ([]Role, error) {
	var list []Role
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.roles.List(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) Permissions(ctx context.Context, p auth.Principal, scope string) ([]Permission, error) {
	var perms []Permission
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sc, err := ParseScope(scope, p.TenantID)
		if err != nil {
			return err
		}
		g, err := s.resolver.ForUser(ctx, tx, p.TenantID, p.UserID, sc)
		perms = g.Permissions.List()
		return err
	})
	return perms, err
}

func (s Service) Require(ctx context.Context, tx pgx.Tx, p auth.Principal, sc Scope, perm Permission) error {
	g, err := s.resolver.ForUser(ctx, tx, p.TenantID, p.UserID, sc)
	if err != nil {
		return err
	}
	if !g.Permissions.Has(perm) {
		return ErrForbidden
	}
	return nil
}

func (s Service) Preview(ctx context.Context, p auth.Principal, roleID, scope string) ([]Permission, error) {
	var perms []Permission
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sc, err := ParseScope(scope, p.TenantID)
		if err != nil {
			return err
		}
		if _, err := s.manager(ctx, tx, p, TenantScope(p.TenantID)); err != nil {
			return err
		}
		if _, err := s.roles.Lock(ctx, tx, roleID); err != nil {
			return err
		}
		g, err := s.resolver.ForRole(ctx, tx, p.TenantID, roleID, sc)
		perms = g.Permissions.List()
		return err
	})
	return perms, err
}

func (s Service) CreateRole(ctx context.Context, p auth.Principal, in RoleInput) (Role, error) {
	var role Role
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.catalog.Validate(in.Permissions); err != nil {
			return err
		}
		if _, err := s.roles.Lock(ctx, tx, ""); err != nil {
			return err
		}
		actor, err := s.manager(ctx, tx, p, TenantScope(p.TenantID))
		if err != nil {
			return err
		}
		if actor.Top == 0 {
			return ErrHierarchy
		}
		if !actor.Permissions.Covers(in.Permissions) {
			return ErrPermissionNotHeld
		}
		role, err = s.roles.Insert(ctx, tx, in, actor.Top)
		return err
	})
	return role, err
}

func (s Service) UpdateRole(ctx context.Context, p auth.Principal, id string, in RoleInput) (Role, error) {
	var role Role
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.catalog.Validate(in.Permissions); err != nil {
			return err
		}
		current, actor, err := s.manageable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if !actor.Permissions.Covers(changed(current.Permissions, in.Permissions)) {
			return ErrPermissionNotHeld
		}
		role, err = s.roles.Update(ctx, tx, id, in)
		return err
	})
	return role, err
}

func (s Service) DeleteRole(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		current, _, err := s.manageable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if current.Default {
			return ErrDefaultRole
		}
		return s.roles.Delete(ctx, tx, current)
	})
}

func (s Service) MoveRole(ctx context.Context, p auth.Principal, id string, position int) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		current, actor, err := s.manageable(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if current.Default {
			return ErrDefaultRole
		}
		if position < 1 || position >= actor.Top {
			return ErrInvalidPosition
		}
		return s.roles.Move(ctx, tx, current, position)
	})
}

func (s Service) Assign(ctx context.Context, p auth.Principal, in AssignInput) (Assignment, error) {
	var a Assignment
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sc, err := ParseScope(in.Scope, p.TenantID)
		if err != nil {
			return err
		}
		if err := s.checkAssignable(ctx, tx, p, in.RoleID, sc); err != nil {
			return err
		}
		a, err = s.assignments.Upsert(ctx, tx, in, sc, &p.UserID)
		return err
	})
	return a, err
}

func (s Service) Unassign(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, sc, err := s.assignments.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.checkAssignable(ctx, tx, p, a.RoleID, sc); err != nil {
			return err
		}
		return s.assignments.Delete(ctx, tx, id)
	})
}

func (s Service) UserRoles(ctx context.Context, p auth.Principal, userID string) ([]Assignment, error) {
	var list []Assignment
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.assignments.ByUser(ctx, tx, userID)
		return err
	})
	return list, err
}

func (s Service) Overrides(ctx context.Context, p auth.Principal, scope string) ([]Override, error) {
	var list []Override
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sc, err := ParseScope(scope, p.TenantID)
		if err != nil {
			return err
		}
		list, err = s.overrides.At(ctx, tx, sc)
		return err
	})
	return list, err
}

func (s Service) SetOverride(ctx context.Context, p auth.Principal, scope string, o Override) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sc, err := ParseScope(scope, p.TenantID)
		if err != nil {
			return err
		}
		if sc.IsTenant() {
			return ErrInvalidScope
		}
		touched := append(slices.Clone(o.Allow), o.Deny...)
		if err := s.catalog.Validate(touched); err != nil {
			return err
		}
		if slices.Contains(touched, Administrator) {
			return ErrUnknownPermission
		}
		role, err := s.roles.Lock(ctx, tx, o.RoleID)
		if err != nil {
			return err
		}
		actor, err := s.manager(ctx, tx, p, sc)
		if err != nil {
			return err
		}
		if role.Position >= actor.Top {
			return ErrHierarchy
		}
		if !actor.Permissions.Covers(touched) {
			return ErrPermissionNotHeld
		}
		return s.overrides.Put(ctx, tx, sc, o)
	})
}

func (s Service) manager(ctx context.Context, tx pgx.Tx, p auth.Principal, sc Scope) (Grant, error) {
	g, err := s.resolver.ForUser(ctx, tx, p.TenantID, p.UserID, sc)
	if err != nil {
		return Grant{}, err
	}
	if !g.Permissions.Has(ManageRoles) {
		return Grant{}, ErrForbidden
	}
	return g, nil
}

func (s Service) manageable(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Role, Grant, error) {
	role, err := s.roles.Lock(ctx, tx, id)
	if err != nil {
		return Role{}, Grant{}, err
	}
	actor, err := s.manager(ctx, tx, p, TenantScope(p.TenantID))
	if err != nil {
		return Role{}, Grant{}, err
	}
	if role.Position >= actor.Top {
		return Role{}, Grant{}, ErrHierarchy
	}
	return role, actor, nil
}

func (s Service) checkAssignable(ctx context.Context, tx pgx.Tx, p auth.Principal, roleID string, sc Scope) error {
	role, err := s.roles.Lock(ctx, tx, roleID)
	if err != nil {
		return err
	}
	if role.Default {
		return ErrDefaultRole
	}
	actor, err := s.manager(ctx, tx, p, sc)
	if err != nil {
		return err
	}
	if role.Position >= actor.Top {
		return ErrHierarchy
	}
	return nil
}

func changed(before, after []Permission) []Permission {
	a, b := NewSet(before...), NewSet(after...)
	var diff []Permission
	for p := range a {
		if _, ok := b[p]; !ok {
			diff = append(diff, p)
		}
	}
	for p := range b {
		if _, ok := a[p]; !ok {
			diff = append(diff, p)
		}
	}
	return diff
}
