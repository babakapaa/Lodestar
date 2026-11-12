package access

import (
	"fmt"
	"regexp"
)

var permissionPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`)

type Catalog struct {
	list  []Permission
	known map[Permission]struct{}
}

func NewCatalog(extra ...string) (Catalog, error) {
	c := Catalog{known: map[Permission]struct{}{}}
	all := corePermissions()
	for _, e := range extra {
		all = append(all, Permission(e))
	}
	for _, p := range all {
		if _, dup := c.known[p]; dup || !permissionPattern.MatchString(string(p)) {
			return Catalog{}, fmt.Errorf("%w: %q", ErrBadPermission, p)
		}
		c.known[p] = struct{}{}
		c.list = append(c.list, p)
	}
	return c, nil
}

func (c Catalog) All() []Permission {
	return append([]Permission(nil), c.list...)
}

func (c Catalog) Validate(perms []Permission) error {
	for _, p := range perms {
		if _, ok := c.known[p]; !ok {
			return fmt.Errorf("%w: %s", ErrUnknownPermission, p)
		}
	}
	return nil
}
