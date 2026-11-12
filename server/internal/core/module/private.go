package module

import "github.com/babakapaa/lodestar/server/internal/core/privacy"

type Private interface {
	PrivacySections() []privacy.Section
}
