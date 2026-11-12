package privacy

import (
	"fmt"
	"regexp"
)

var sectionPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$`)

type Sections []Section

func NewSections(list ...Section) (Sections, error) {
	seen := map[string]bool{}
	for _, s := range list {
		if !sectionPattern.MatchString(s.Name()) || seen[s.Name()] {
			return nil, fmt.Errorf("%w: %q", ErrBadSection, s.Name())
		}
		seen[s.Name()] = true
	}
	return Sections(list), nil
}
