package access

import (
	"regexp"
	"strings"
)

const TenantKind = "tenant"

var scopeKindPattern = regexp.MustCompile(`^[a-z][a-z0-9_.]{0,39}$`)

type Scope struct {
	Kind string
	ID   string
}

func TenantScope(tenantID string) Scope {
	return Scope{Kind: TenantKind, ID: tenantID}
}

func ParseScope(value, tenantID string) (Scope, error) {
	if value == "" || value == TenantKind {
		return TenantScope(tenantID), nil
	}
	kind, id, ok := strings.Cut(value, ":")
	if !ok || kind == TenantKind || !scopeKindPattern.MatchString(kind) || id == "" || len(id) > 100 {
		return Scope{}, ErrInvalidScope
	}
	return Scope{Kind: kind, ID: id}, nil
}

func (s Scope) IsTenant() bool {
	return s.Kind == TenantKind
}

func (s Scope) String() string {
	if s.IsTenant() {
		return TenantKind
	}
	return s.Kind + ":" + s.ID
}
