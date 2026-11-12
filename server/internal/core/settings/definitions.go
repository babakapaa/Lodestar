package settings

import "github.com/babakapaa/lodestar/server/internal/core/schema"

func CoreSettings() []schema.Definition {
	return []schema.Definition{
		{Key: "core.default_language", Type: schema.Choice, Options: []string{"system", "kk", "ru", "en"}, Default: "system"},
		{Key: "core.default_theme", Type: schema.Choice, Options: []string{"paper", "graphite", "lab", "campus"}, Default: "graphite"},
	}
}
