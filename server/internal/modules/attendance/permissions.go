package attendance

import "github.com/babakapaa/lodestar/server/internal/core/access"

const Manage access.Permission = "attendance.manage"

func (Module) Permissions() []string {
	return []string{string(Manage)}
}
