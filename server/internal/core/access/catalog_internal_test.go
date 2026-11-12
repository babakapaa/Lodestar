package access

import (
	"errors"
	"testing"
)

func TestCatalogRejectsBadKeys(t *testing.T) {
	for _, key := range []string{"nodot", "Bad.Key", string(ManageRoles)} {
		if _, err := NewCatalog(key); !errors.Is(err, ErrBadPermission) {
			t.Errorf("NewCatalog(%q) err = %v", key, err)
		}
	}
}

func TestParseScope(t *testing.T) {
	s, err := ParseScope("group:g1", "t")
	if err != nil || s.Kind != "group" || s.ID != "g1" {
		t.Fatalf("ParseScope = %+v, %v", s, err)
	}
	if s, _ := ParseScope("", "t"); !s.IsTenant() || s.ID != "t" {
		t.Fatalf("empty scope = %+v", s)
	}
	for _, bad := range []string{"group", "tenant:x", "Group:1", "group:"} {
		if _, err := ParseScope(bad, "t"); !errors.Is(err, ErrInvalidScope) {
			t.Errorf("ParseScope(%q) err = %v", bad, err)
		}
	}
}
