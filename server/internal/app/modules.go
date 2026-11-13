package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
	"github.com/babakapaa/lodestar/server/internal/modules/league"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
	"github.com/babakapaa/lodestar/server/internal/modules/notes"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
	"github.com/babakapaa/lodestar/server/internal/modules/team"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func Modules() []module.Module {
	return []module.Module{study.Module{}, notes.Module{}, navigator.Module{}, chat.Module{}, team.Module{}, league.Module{}}
}

func CoreMigrations() []database.Migrations {
	return []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), privacy.Migrations(), auth.Migrations(), access.Migrations(), structure.Migrations(),
		profile.Migrations(), settings.Migrations(), live.Migrations(), notify.Migrations(),
		files.Migrations(), delivery.Migrations(), lti.Migrations(),
	}
}
