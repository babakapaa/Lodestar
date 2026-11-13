package module

type LinkSources []LinkSource

func (l LinkSources) Find(kind string) (LinkSource, bool) {
	for _, source := range l {
		if source.Kind() == kind {
			return source, true
		}
	}
	return nil, false
}
