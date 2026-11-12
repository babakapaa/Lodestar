package notify

import (
	"fmt"
	"regexp"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

var typePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`)

type Catalog struct {
	list  []module.NotificationType
	byKey map[string]module.NotificationType
}

func NewCatalog(types ...module.NotificationType) (Catalog, error) {
	c := Catalog{byKey: map[string]module.NotificationType{}}
	for _, t := range types {
		if _, dup := c.byKey[t.Key]; dup || !typePattern.MatchString(t.Key) || len(t.Texts) == 0 || len(t.Names) == 0 {
			return Catalog{}, fmt.Errorf("%w: %q", ErrBadType, t.Key)
		}
		if t.Required {
			t.InApp = true
		}
		c.byKey[t.Key] = t
		c.list = append(c.list, t)
	}
	return c, nil
}

func (c Catalog) Lookup(key string) (module.NotificationType, error) {
	t, ok := c.byKey[key]
	if !ok {
		return module.NotificationType{}, fmt.Errorf("%w: %s", ErrUnknownType, key)
	}
	return t, nil
}

func (c Catalog) All() []module.NotificationType {
	return append([]module.NotificationType(nil), c.list...)
}
