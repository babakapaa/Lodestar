package module

import (
	"fmt"
	"strings"
)

func checkSearchKinds(m Module) error {
	s, ok := m.(Searchable)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, src := range s.SearchSources() {
		if !strings.HasPrefix(src.Kind(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignSearchKind, src.Kind(), name)
		}
	}
	return nil
}
