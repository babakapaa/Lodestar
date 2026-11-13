package navigator

import "github.com/babakapaa/lodestar/server/internal/core/schema"

const JournalPreference = "navigator.journal"

func (Module) Preferences() []schema.Definition {
	return []schema.Definition{{Key: JournalPreference, Type: schema.Bool, Default: true}}
}
