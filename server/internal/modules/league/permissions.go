package league

import "github.com/babakapaa/lodestar/server/internal/core/access"

const ManageAwards access.Permission = "league.awards.manage"

func (Module) Permissions() []string {
	return []string{string(ManageAwards)}
}
