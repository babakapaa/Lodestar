package app

import (
	"slices"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/stats"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

func CorePermissions() []string {
	return slices.Concat(structure.Permissions(), settings.Permissions(), audit.Permissions(), lti.Permissions(), stats.Permissions())
}
