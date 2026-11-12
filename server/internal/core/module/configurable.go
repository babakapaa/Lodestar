package module

import "github.com/babakapaa/lodestar/server/internal/core/schema"

type Configurable interface {
	Settings() []schema.Definition
}
