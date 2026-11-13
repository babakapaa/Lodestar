package lti

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type ContextView struct {
	ID         string              `json:"id"`
	Title      string              `json:"title"`
	Label      string              `json:"label"`
	Platform   string              `json:"platform"`
	Instructor bool                `json:"instructor"`
	Link       *module.LinkTarget  `json:"link"`
	Options    []module.LinkTarget `json:"options"`
	Source     string              `json:"assignmentSource"`
	CanChoose  bool                `json:"sourceEditable"`
	CanSync    bool                `json:"canSync"`
	SyncedAt   *time.Time          `json:"syncedAt"`
	SyncError  *string             `json:"syncError"`
}
