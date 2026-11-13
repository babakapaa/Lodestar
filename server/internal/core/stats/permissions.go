package stats

import "github.com/babakapaa/lodestar/server/internal/core/access"

const (
	View access.Permission = "core.stats.view"
	Full access.Permission = "core.stats.full"
)

func Permissions() []string {
	return []string{string(View), string(Full)}
}
