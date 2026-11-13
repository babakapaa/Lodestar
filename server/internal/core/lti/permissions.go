package lti

import "github.com/babakapaa/lodestar/server/internal/core/access"

const ManageIntegrations access.Permission = "core.integrations.manage"

func Permissions() []string {
	return []string{string(ManageIntegrations)}
}
