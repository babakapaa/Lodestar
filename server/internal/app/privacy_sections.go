package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

func privacySections(ext Integrations) []privacy.Section {
	return []privacy.Section{
		identity.PrivacySection{},
		profile.PrivacySection{},
		auth.PrivacySection{},
		access.PrivacySection{},
		structure.PrivacySection{},
		notify.PrivacySection{},
		delivery.PrivacySection{},
		files.NewPrivacySection(ext.Storage),
		audit.PrivacySection{},
	}
}
