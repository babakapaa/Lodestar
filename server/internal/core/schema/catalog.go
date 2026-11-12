package schema

import (
	"fmt"
	"regexp"
)

var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`)

type Catalog struct {
	list  []Definition
	byKey map[string]Definition
}

func NewCatalog(defs ...Definition) (Catalog, error) {
	c := Catalog{byKey: map[string]Definition{}}
	for _, d := range defs {
		if _, dup := c.byKey[d.Key]; dup || !keyPattern.MatchString(d.Key) {
			return Catalog{}, fmt.Errorf("%w: %q", ErrBadDefinition, d.Key)
		}
		if err := d.Validate(d.DefaultValue()); err != nil {
			return Catalog{}, fmt.Errorf("%w: default of %s: %w", ErrBadDefinition, d.Key, err)
		}
		c.byKey[d.Key] = d
		c.list = append(c.list, d)
	}
	return c, nil
}

func (c Catalog) All() []Definition {
	return append([]Definition(nil), c.list...)
}

func (c Catalog) Lookup(key string) (Definition, error) {
	d, ok := c.byKey[key]
	if !ok {
		return Definition{}, fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}
	return d, nil
}
