package notes

import (
	"strconv"
	"strings"
)

type crossrefWork struct {
	Type      string   `json:"type"`
	Title     []string `json:"title"`
	Container []string `json:"container-title"`
	Publisher string   `json:"publisher"`
	Place     string   `json:"publisher-location"`
	Volume    string   `json:"volume"`
	Issue     string   `json:"issue"`
	Page      string   `json:"page"`
	DOI       string   `json:"DOI"`
	URL       string   `json:"URL"`
	ISBN      []string `json:"ISBN"`
	Author    []struct {
		Given  string `json:"given"`
		Family string `json:"family"`
		Name   string `json:"name"`
	} `json:"author"`
	Issued struct {
		Parts [][]int `json:"date-parts"`
	} `json:"issued"`
}

func (w crossrefWork) fields() SourceFields {
	f := SourceFields{
		Kind:      crossrefKind(w.Type),
		Title:     first(w.Title),
		Container: first(w.Container),
		Publisher: w.Publisher,
		City:      w.Place,
		Volume:    w.Volume,
		Issue:     w.Issue,
		Pages:     strings.ReplaceAll(w.Page, "-", "–"),
		DOI:       w.DOI,
		URL:       w.URL,
		ISBN:      first(w.ISBN),
	}
	for _, a := range w.Author {
		name := strings.TrimSpace(strings.TrimSpace(a.Family + ", " + a.Given))
		name = strings.Trim(name, ", ")
		if a.Name != "" && name == "" {
			name = a.Name
		}
		if name != "" {
			f.Authors = append(f.Authors, name)
		}
	}
	if len(w.Issued.Parts) > 0 && len(w.Issued.Parts[0]) > 0 && w.Issued.Parts[0][0] > 0 {
		f.Year = strconv.Itoa(w.Issued.Parts[0][0])
	}
	return f
}

func crossrefKind(t string) string {
	switch t {
	case "journal-article", "proceedings-article", "posted-content":
		return "article"
	case "book", "monograph", "edited-book", "reference-book":
		return "book"
	case "book-chapter", "book-section", "book-part":
		return "chapter"
	case "dissertation":
		return "thesis"
	default:
		return "other"
	}
}

func first(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return strings.TrimSpace(list[0])
}
