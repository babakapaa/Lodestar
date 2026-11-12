package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

func CorePermissions() []string {
	return append(append(structure.Permissions(), settings.Permissions()...), audit.Permissions()...)
}
