package app

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/safehttp"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Core struct {
	Access        access.Service
	Settings      settings.Service
	Notifications notify.Catalog
	Kernel        module.Kernel
}

func NewCore(pool *pgxpool.Pool, registry *module.Registry, box secretbox.Box) (Core, error) {
	tx := database.NewTenantTx(pool)
	catalog, err := access.NewCatalog(append(CorePermissions(), registry.Permissions()...)...)
	if err != nil {
		return Core{}, err
	}
	tenantSettings, err := schema.NewCatalog(append(coreSettings(), registry.Settings()...)...)
	if err != nil {
		return Core{}, err
	}
	notifications, err := notify.NewCatalog(registry.NotificationTypes()...)
	if err != nil {
		return Core{}, err
	}
	accessService := access.NewService(tx, catalog, structure.Tree{})
	settingsService := settings.NewService(tx, accessService, tenantSettings, registry)
	return Core{
		Access:        accessService,
		Settings:      settingsService,
		Notifications: notifications,
		Kernel: module.Kernel{
			Pool:     pool,
			Tenants:  tx,
			Access:   accessService,
			Notifier: moduleNotifier{notifier: notify.NewNotifier(notifications)},
			Settings: settingsService,
			Secrets:  box,
			Outbound: safehttp.NewClient(15 * time.Second),
		},
	}, nil
}
