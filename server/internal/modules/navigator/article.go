package navigator

import (
	"strconv"
	"strings"
)

type Article struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Year    int      `json:"year"`
	Venue   string   `json:"venue"`
	DOI     string   `json:"doi"`
	URL     string   `json:"url"`
}

func (a Article) label() string {
	parts := []string{}
	if len(a.Authors) > 0 {
		parts = append(parts, a.Authors[0])
	}
	parts = append(parts, a.Title)
	if a.Year > 0 {
		parts = append(parts, strconv.Itoa(a.Year))
	}
	label := strings.Join(parts, ". ")
	if a.DOI != "" {
		label += ". DOI: " + a.DOI
	}
	return label
}
