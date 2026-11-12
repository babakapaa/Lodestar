package meta

import "github.com/babakapaa/lodestar/server/internal/core/module"

type Response struct {
	Version string            `json:"version"`
	Modules []module.Manifest `json:"modules"`
}
