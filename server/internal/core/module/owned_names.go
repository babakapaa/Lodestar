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

func checkPrivacySections(m Module) error {
	s, ok := m.(Private)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, section := range s.PrivacySections(Kernel{}) {
		if !strings.HasPrefix(section.Name(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignSection, section.Name(), name)
		}
	}
	return nil
}

func checkRosterKinds(m Module) error {
	s, ok := m.(Rostering)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, src := range s.RosterSources(Kernel{}) {
		if !strings.HasPrefix(src.Kind(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignRoster, src.Kind(), name)
		}
	}
	return nil
}

func checkTeamScopeKinds(m Module) error {
	s, ok := m.(Teaming)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, src := range s.TeamScopeSources(Kernel{}) {
		if !strings.HasPrefix(src.Kind(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignTeamScope, src.Kind(), name)
		}
	}
	return nil
}

func checkScoreKinds(m Module) error {
	s, ok := m.(Scoring)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, src := range s.ScoreSources(Kernel{}) {
		if !strings.HasPrefix(src.Kind(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignScoreSource, src.Kind(), name)
		}
	}
	return nil
}

func checkLinkKinds(m Module) error {
	s, ok := m.(LinkProviding)
	if !ok {
		return nil
	}
	name := m.Manifest().Name
	for _, src := range s.LinkSources(Kernel{}) {
		if !strings.HasPrefix(src.Kind(), name+".") {
			return fmt.Errorf("%w: %s in module %s", ErrForeignLinkSource, src.Kind(), name)
		}
	}
	return nil
}
