package profile

import "github.com/babakapaa/lodestar/server/internal/core/schema"

func CorePreferences() []schema.Definition {
	minScale, maxScale := schema.Bounds(80, 200)
	return []schema.Definition{
		{Key: "core.language", Type: schema.Choice, Options: []string{"system", "kk", "ru", "en"}, Default: "system"},
		{Key: "core.theme", Type: schema.Choice, Options: []string{"tenant", "paper", "graphite", "lab", "campus"}, Default: "tenant"},
		{Key: "core.appearance", Type: schema.Choice, Options: []string{"system", "light", "dark"}, Default: "system"},
		{Key: "core.text_scale", Type: schema.Integer, Min: minScale, Max: maxScale, Default: 100},
		{Key: "core.reduced_motion", Type: schema.Bool, Default: false},
	}
}
