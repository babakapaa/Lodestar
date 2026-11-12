package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func Modules() []module.Module {
	return []module.Module{}
}

func CoreMigrations() []database.Migrations {
	return []database.Migrations{tenant.Migrations()}
}
