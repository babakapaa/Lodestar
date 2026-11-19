package navigator

import "github.com/babakapaa/lodestar/server/internal/core/schema"

const (
	SupportSetting = "navigator.support_contacts"
	MemorySetting  = "navigator.memory_messages"
	LibrarySetting = "navigator.open_library"
)

func (Module) Settings() []schema.Definition {
	memoryMin, memoryMax := 0, 200
	return []schema.Definition{
		{Key: SupportSetting, Type: schema.Text, MaxLength: 2000, Default: ""},
		{Key: MemorySetting, Type: schema.Integer, Default: 30, Min: &memoryMin, Max: &memoryMax},
		{Key: LibrarySetting, Type: schema.Bool, Default: true},
	}
}
