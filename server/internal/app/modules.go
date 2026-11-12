package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func Modules() []module.Module {
	return []module.Module{}
}

func CoreMigrations() []database.Migrations {
	return []database.Migrations{tenant.Migrations(), identity.Migrations(), auth.Migrations(), access.Migrations()}
}
