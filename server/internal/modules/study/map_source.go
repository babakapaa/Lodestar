package study

import (
	"slices"
	"strings"
)

var scholarlyKinds = []string{"article", "chapter", "thesis"}

type MapSource struct {
	SourceID  string `json:"sourceId"`
	Reference string `json:"reference"`
	Kind      string `json:"kind"`
	DOI       string `json:"doi"`
	URL       string `json:"url"`
	Found     string `json:"found"`
	Why       string `json:"why"`
	Taken     string `json:"taken"`
}

func (s MapSource) Scholarly() bool {
	return s.DOI != "" || slices.Contains(scholarlyKinds, s.Kind)
}

func (s MapSource) normalized() (MapSource, bool) {
	for _, field := range []*string{&s.SourceID, &s.Reference, &s.Kind, &s.DOI, &s.URL, &s.Found, &s.Why, &s.Taken} {
		*field = strings.TrimSpace(*field)
		if len(*field) > maxMapText {
			return s, false
		}
	}
	return s, s.Reference != ""
}
