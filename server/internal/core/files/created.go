package files

import "github.com/babakapaa/lodestar/server/internal/platform/objectstore"

type Created struct {
	File   File               `json:"file"`
	Upload objectstore.Upload `json:"upload"`
}
