package study

import "strings"

const (
	maxMapSources = 200
	maxMapQueries = 100
	maxMapText    = 2000
	maxMapPath    = 20000
)

const (
	MapMissingSources   = "sources"
	MapMissingScholarly = "scholarly"
	MapMissingWhy       = "why"
	MapMissingPath      = "path"
)

type SearchMap struct {
	Used     []MapSource `json:"used"`
	Rejected []MapSource `json:"rejected"`
	Path     string      `json:"path"`
	Queries  []string    `json:"queries"`
}

func (m SearchMap) normalized() (SearchMap, error) {
	if len(m.Used)+len(m.Rejected) > maxMapSources || len(m.Queries) > maxMapQueries {
		return m, ErrInvalidSearchMap
	}
	out := SearchMap{Path: strings.TrimSpace(m.Path), Used: []MapSource{}, Rejected: []MapSource{}, Queries: []string{}}
	if len(out.Path) > maxMapPath {
		return m, ErrInvalidSearchMap
	}
	for _, list := range []struct {
		in  []MapSource
		out *[]MapSource
	}{{m.Used, &out.Used}, {m.Rejected, &out.Rejected}} {
		for _, s := range list.in {
			s, ok := s.normalized()
			if !ok {
				return m, ErrInvalidSearchMap
			}
			*list.out = append(*list.out, s)
		}
	}
	for _, q := range m.Queries {
		q = strings.TrimSpace(q)
		if len(q) > maxMapText {
			return m, ErrInvalidSearchMap
		}
		if q != "" {
			out.Queries = append(out.Queries, q)
		}
	}
	return out, nil
}

func (m SearchMap) Missing(r MapRules) []string {
	missing := []string{}
	if len(m.Used) == 0 || len(m.Used) < r.MinSources {
		missing = append(missing, MapMissingSources)
	}
	scholarly := 0
	for _, s := range m.Used {
		if s.Scholarly() {
			scholarly++
		}
	}
	if len(m.Used) > 0 && scholarly*100 < r.MinScholarly*len(m.Used) {
		missing = append(missing, MapMissingScholarly)
	}
	for _, s := range append(append([]MapSource{}, m.Used...), m.Rejected...) {
		if s.Why == "" {
			missing = append(missing, MapMissingWhy)
			break
		}
	}
	if m.Path == "" {
		missing = append(missing, MapMissingPath)
	}
	return missing
}
