package notes

import (
	"strings"
	"unicode/utf8"
)

const (
	maxSourceText    = 1000
	maxSourceNote    = 10000
	maxSourceAuthors = 100
)

var sourceKinds = map[string]bool{
	"book": true, "article": true, "chapter": true, "web": true, "thesis": true, "other": true,
}

type SourceFields struct {
	Kind       string   `json:"kind"`
	Title      string   `json:"title"`
	Authors    []string `json:"authors"`
	Year       string   `json:"year"`
	Container  string   `json:"container"`
	Publisher  string   `json:"publisher"`
	City       string   `json:"city"`
	Volume     string   `json:"volume"`
	Issue      string   `json:"issue"`
	Pages      string   `json:"pages"`
	URL        string   `json:"url"`
	DOI        string   `json:"doi"`
	ISBN       string   `json:"isbn"`
	AccessedAt string   `json:"accessedAt"`
	Note       string   `json:"note"`
}

func (f SourceFields) normalized() (SourceFields, error) {
	f.Kind = strings.TrimSpace(f.Kind)
	if f.Kind == "" {
		f.Kind = "other"
	}
	if !sourceKinds[f.Kind] || len(f.Authors) > maxSourceAuthors || utf8.RuneCountInString(f.Note) > maxSourceNote {
		return f, ErrInvalidSource
	}
	for _, field := range []*string{
		&f.Title, &f.Year, &f.Container, &f.Publisher, &f.City, &f.Volume, &f.Issue,
		&f.Pages, &f.URL, &f.DOI, &f.ISBN, &f.AccessedAt,
	} {
		*field = strings.TrimSpace(*field)
		if utf8.RuneCountInString(*field) > maxSourceText {
			return f, ErrInvalidSource
		}
	}
	if f.Title == "" {
		return f, ErrInvalidSource
	}
	authors := make([]string, 0, len(f.Authors))
	for _, a := range f.Authors {
		a = strings.TrimSpace(a)
		if utf8.RuneCountInString(a) > maxSourceText {
			return f, ErrInvalidSource
		}
		if a != "" {
			authors = append(authors, a)
		}
	}
	f.Authors = authors
	f.DOI = NormalizeDOI(f.DOI)
	return f, nil
}

func NormalizeDOI(doi string) string {
	doi = strings.TrimSpace(doi)
	lower := strings.ToLower(doi)
	for _, prefix := range []string{"https://doi.org/", "http://doi.org/", "https://dx.doi.org/", "http://dx.doi.org/", "doi:"} {
		if strings.HasPrefix(lower, prefix) {
			doi = strings.TrimSpace(doi[len(prefix):])
			break
		}
	}
	return doi
}
