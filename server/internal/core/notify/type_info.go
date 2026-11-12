package notify

import "github.com/babakapaa/lodestar/server/internal/core/module"

type TypeInfo struct {
	Key      string                             `json:"key"`
	Required bool                               `json:"required"`
	Names    map[string]string                  `json:"names"`
	Texts    map[string]module.NotificationText `json:"texts"`
}
