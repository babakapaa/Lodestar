package module

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Kernel struct {
	Pool     *pgxpool.Pool
	Tenants  database.TenantTx
	Access   access.Service
	Notifier Notifier
	Settings Settings
	Secrets  secretbox.Box
	Outbound *http.Client
	Model    llm.Model
	Briefs   Briefs
}
