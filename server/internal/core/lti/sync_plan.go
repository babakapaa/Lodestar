package lti

import "github.com/babakapaa/lodestar/server/internal/core/module"

type syncPlan struct {
	tenantID string
	context  Context
	platform Platform
	key      Key
	source   module.LinkSource
	mode     string
	actor    string
}

func (p syncPlan) has(scope string) bool {
	for _, s := range p.context.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}
