package module

type Rosters []RosterSource

func (r Rosters) Find(kind string) (RosterSource, bool) {
	for _, source := range r {
		if source.Kind() == kind {
			return source, true
		}
	}
	return nil, false
}
