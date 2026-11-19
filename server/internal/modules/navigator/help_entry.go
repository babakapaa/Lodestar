package navigator

import "github.com/babakapaa/lodestar/server/internal/core/module"

type HelpEntry struct {
	ScopeKind string             `json:"scopeKind"`
	ScopeID   string             `json:"scopeId"`
	Summary   module.HelpSummary `json:"summary"`
}
