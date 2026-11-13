package lti

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type DeepLink struct {
	ID             string              `json:"id"`
	PlatformID     string              `json:"-"`
	UserID         string              `json:"-"`
	DeploymentID   string              `json:"-"`
	ReturnURL      string              `json:"-"`
	Data           string              `json:"-"`
	Platform       string              `json:"platform"`
	AcceptMultiple bool                `json:"acceptMultiple"`
	Targets        []module.LinkTarget `json:"targets"`
	CreatedAt      time.Time           `json:"createdAt"`
}
