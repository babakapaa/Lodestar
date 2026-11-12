package module

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Kernel struct {
	Pool     *pgxpool.Pool
	Tenants  database.TenantTx
	Access   access.Service
	Notifier Notifier
	Settings Settings
}
