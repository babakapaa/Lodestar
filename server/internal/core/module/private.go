package module

import "github.com/babakapaa/lodestar/server/internal/core/privacy"

type Private interface {
	PrivacySections(k Kernel) []privacy.Section
}
