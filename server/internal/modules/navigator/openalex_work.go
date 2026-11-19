package navigator

import "strings"

type openAlexWork struct {
	Title   string `json:"display_name"`
	DOI     string `json:"doi"`
	Year    int    `json:"publication_year"`
	Authors []struct {
		Author struct {
			Name string `json:"display_name"`
		} `json:"author"`
	} `json:"authorships"`
	Location struct {
		Source *struct {
			Name string `json:"display_name"`
		} `json:"source"`
		Landing string `json:"landing_page_url"`
	} `json:"primary_location"`
}

func (w openAlexWork) article() (Article, bool) {
	a := Article{Title: strings.Join(strings.Fields(w.Title), " "), Year: w.Year, URL: w.Location.Landing}
	if a.Title == "" {
		return Article{}, false
	}
	if doi := strings.TrimPrefix(strings.TrimSpace(w.DOI), "https://doi.org/"); strings.HasPrefix(doi, "10.") {
		a.DOI, a.URL = doi, "https://doi.org/"+doi
	}
	if w.Location.Source != nil {
		a.Venue = w.Location.Source.Name
	}
	for i, au := range w.Authors {
		if i == maxArticleAuthors {
			break
		}
		if name := strings.TrimSpace(au.Author.Name); name != "" {
			a.Authors = append(a.Authors, name)
		}
	}
	if a.Authors == nil {
		a.Authors = []string{}
	}
	return a, true
}
