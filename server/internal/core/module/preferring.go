package module

import "github.com/babakapaa/lodestar/server/internal/core/schema"

type Preferring interface {
	Preferences() []schema.Definition
}
