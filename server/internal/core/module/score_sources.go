package module

type ScoreSources []ScoreSource

func (s ScoreSources) Find(kind string) (ScoreSource, bool) {
	for _, source := range s {
		if source.Kind() == kind {
			return source, true
		}
	}
	return nil, false
}
