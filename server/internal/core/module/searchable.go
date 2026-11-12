package module

import "github.com/babakapaa/lodestar/server/internal/core/search"

type Searchable interface {
	SearchSources() []search.Source
}
