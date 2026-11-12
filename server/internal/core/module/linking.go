package module

import "github.com/babakapaa/lodestar/server/internal/core/files"

type Linking interface {
	FileLinks(k Kernel) []files.Link
}
