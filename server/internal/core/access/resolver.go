package access

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Resolver struct {
	tree        ScopeTree
	catalog     Catalog
	roles       RoleStore
	assignments AssignmentStore
	overrides   OverrideStore
}

func NewResolver(tree ScopeTree, catalog Catalog) Resolver {
	return Resolver{tree: tree, catalog: catalog}
}

func (r Resolver) ForUser(ctx context.Context, tx pgx.Tx, tenantID, userID string, s Scope) (Grant, error) {
	chain, err := r.tree.Chain(ctx, tx, tenantID, s)
	if err != nil {
		return Grant{}, err
	}
	roleIDs, err := r.assignments.ActiveRoleIDs(ctx, tx, userID, chain)
	if err != nil {
		return Grant{}, err
	}
	return r.evaluate(ctx, tx, chain, roleIDs)
}

func (r Resolver) ForRole(ctx context.Context, tx pgx.Tx, tenantID, roleID string, s Scope) (Grant, error) {
	chain, err := r.tree.Chain(ctx, tx, tenantID, s)
	if err != nil {
		return Grant{}, err
	}
	return r.evaluate(ctx, tx, chain, []string{roleID})
}

func (r Resolver) evaluate(ctx context.Context, tx pgx.Tx, chain []Scope, roleIDs []string) (Grant, error) {
	roles, err := r.roles.ByIDs(ctx, tx, roleIDs)
	if err != nil {
		return Grant{}, err
	}
	g := Grant{Permissions: Set{}}
	ids := make([]string, 0, len(roles))
	for _, role := range roles {
		g.Permissions.Add(role.Permissions...)
		g.Top = max(g.Top, role.Position)
		ids = append(ids, role.ID)
	}
	if g.Permissions.Has(Administrator) {
		g.Permissions = NewSet(r.catalog.All()...)
		return g, nil
	}
	overrides, err := r.overrides.InChain(ctx, tx, chain, ids)
	if err != nil {
		return Grant{}, err
	}
	for _, s := range chain {
		allow, deny := Set{}, Set{}
		for _, so := range overrides {
			if so.scope == s {
				allow.Add(so.override.Allow...)
				deny.Add(so.override.Deny...)
			}
		}
		g.Permissions.Remove(deny.List()...)
		g.Permissions.Add(allow.List()...)
	}
	return g, nil
}
