package settings

import "github.com/babakapaa/lodestar/server/internal/core/access"

const ManageSettings access.Permission = "core.settings.manage"

func Permissions() []string {
	return []string{string(ManageSettings)}
}
