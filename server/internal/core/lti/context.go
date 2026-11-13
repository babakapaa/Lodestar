package lti

import "time"

type Context struct {
	ID           string
	PlatformID   string
	DeploymentID string
	ContextID    string
	Title        string
	Label        string
	MembersURL   *string
	LineItemsURL *string
	Scopes       []string
	OriginKind   *string
	OriginID     *string
	LinkedBy     *string
	Source       *string
	SyncedAt     *time.Time
	SyncError    *string
}
