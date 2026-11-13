package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
)

func NotificationTypes(registry *module.Registry) []module.NotificationType {
	return append(setup.NotificationTypes(), registry.NotificationTypes()...)
}
