package structure

import "github.com/babakapaa/lodestar/server/internal/core/access"

const (
	ScopeKind                         = "unit"
	ManageStructure access.Permission = "core.structure.manage"
	ManageMembers   access.Permission = "core.structure.members"
)

func Permissions() []string {
	return []string{string(ManageStructure), string(ManageMembers)}
}
