package demo

import "github.com/babakapaa/lodestar/server/internal/core/tenant"

type Result struct {
	Tenant   tenant.Tenant
	Callback string
	People   map[string]string
}
