package league

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

type SolutionWork struct {
	Text  string   `json:"text"`
	Links []string `json:"links"`
}

func (w SolutionWork) normalized() (SolutionWork, error) {
	w.Text = strings.TrimSpace(w.Text)
	if utf8.RuneCountInString(w.Text) > 20000 || len(w.Links) > 10 || (w.Text == "" && len(w.Links) == 0) {
		return w, ErrInvalidSolution
	}
	links := make([]string, 0, len(w.Links))
	for _, link := range w.Links {
		link = strings.TrimSpace(link)
		u, err := url.Parse(link)
		if err != nil || len(link) > 500 || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
			return w, ErrInvalidSolution
		}
		links = append(links, link)
	}
	w.Links = links
	return w, nil
}
