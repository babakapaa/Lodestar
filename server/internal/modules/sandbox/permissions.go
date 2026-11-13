package sandbox

import "github.com/babakapaa/lodestar/server/internal/core/access"

const Manage access.Permission = "sandbox.manage"

func (Module) Permissions() []string {
	return []string{string(Manage)}
}
