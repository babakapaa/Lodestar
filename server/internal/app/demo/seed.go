package demo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const (
	TenantSlug  = "demo"
	TenantName  = "Демо-университет"
	ProviderKey = "demo"
	TeacherRole = "Преподаватель"
)

func Seed(ctx context.Context, pool *pgxpool.Pool, box secretbox.Box, registry *module.Registry, cfg Config) (Result, error) {
	core, err := app.NewCore(pool, registry, box)
	if err != nil {
		return Result{}, err
	}
	t, err := ensureTenant(ctx, pool)
	if err != nil {
		return Result{}, err
	}
	tx := database.NewTenantTx(pool)
	_, err = auth.NewProviderStore(tx, box).Put(ctx, auth.Provider{
		TenantID:     t.ID,
		Key:          ProviderKey,
		DisplayName:  "Демо-вход",
		Issuer:       cfg.Issuer,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"openid", "email", "profile"},
	})
	if err != nil {
		return Result{}, err
	}
	ids := map[string]string{}
	users := identity.NewStore(tx)
	for _, person := range People() {
		u, err := users.UpsertExternal(ctx, t.ID, identity.ExternalIdentity{
			Issuer:      cfg.Issuer,
			Subject:     person.Subject,
			Email:       person.Email,
			DisplayName: person.Name,
		})
		if err != nil {
			return Result{}, err
		}
		ids[person.Subject] = u.ID
	}
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, t.ID, ids[Admin.Subject]); err != nil {
		return Result{}, err
	}
	if err := grantTeacher(ctx, core.Access, t.ID, ids[Admin.Subject], ids[Teacher.Subject]); err != nil {
		return Result{}, err
	}
	if err := seedStudy(ctx, core, t.ID, ids, time.Now()); err != nil {
		return Result{}, err
	}
	return Result{Tenant: t, Callback: auth.CallbackURL(cfg.PublicURL, t.Slug, ProviderKey), People: ids}, nil
}

func ensureTenant(ctx context.Context, pool *pgxpool.Pool) (tenant.Tenant, error) {
	store := tenant.NewStore(pool)
	t, err := store.BySlug(ctx, TenantSlug)
	if errors.Is(err, tenant.ErrNotFound) {
		t, err = store.Create(ctx, TenantSlug, TenantName)
	}
	if err != nil {
		return tenant.Tenant{}, err
	}
	return t, access.NewBootstrap(database.NewTenantTx(pool)).EnsureDefaults(ctx, t.ID)
}

func grantTeacher(ctx context.Context, service access.Service, tenantID, adminID, teacherID string) error {
	admin := auth.Principal{TenantID: tenantID, UserID: adminID}
	roles, err := service.Roles(ctx, admin)
	if err != nil {
		return err
	}
	roleID := ""
	for _, r := range roles {
		if r.Name == TeacherRole {
			roleID = r.ID
		}
	}
	if roleID == "" {
		created, err := service.CreateRole(ctx, admin, access.RoleInput{
			Name:        TeacherRole,
			Color:       "#2563EB",
			Permissions: []access.Permission{structure.ManageMembers},
		})
		if err != nil {
			return err
		}
		roleID = created.ID
	}
	_, err = service.Assign(ctx, admin, access.AssignInput{RoleID: roleID, UserID: teacherID})
	return err
}
