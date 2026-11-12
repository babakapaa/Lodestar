package module

import (
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Kernel struct {
	Tenants  database.TenantTx
	Access   access.Service
	Notifier Notifier
}
