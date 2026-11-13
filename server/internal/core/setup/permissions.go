package setup

import "github.com/babakapaa/lodestar/server/internal/core/access"

const (
	Manage access.Permission = "core.setup.manage"
	Legal  access.Permission = "core.setup.legal"
	Sign   access.Permission = "core.setup.sign"
)

func Permissions() []string {
	return []string{string(Manage), string(Legal), string(Sign)}
}
