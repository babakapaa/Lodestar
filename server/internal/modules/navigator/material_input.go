package navigator

import (
	"net/url"
	"strings"
	"unicode/utf8"
)

const (
	maxMaterialTitle = 300
	maxMaterialLink  = 2000
	maxMaterialText  = 1500000
)

type MaterialInput struct {
	Title     string `json:"title"`
	Reference string `json:"reference"`
	URL       string `json:"url"`
	Text      string `json:"text"`
}

func (in MaterialInput) normalized() (MaterialInput, error) {
	in.Title = strings.Join(strings.Fields(in.Title), " ")
	in.Reference = strings.Join(strings.Fields(in.Reference), " ")
	in.URL = strings.TrimSpace(in.URL)
	in.Text = strings.TrimSpace(strings.ReplaceAll(in.Text, "\r\n", "\n"))
	if in.Title == "" || utf8.RuneCountInString(in.Title) > maxMaterialTitle ||
		utf8.RuneCountInString(in.Reference) > maxMaterialLink || len(in.URL) > maxMaterialLink ||
		in.Text == "" || utf8.RuneCountInString(in.Text) > maxMaterialText {
		return in, ErrInvalidMaterial
	}
	if in.URL != "" {
		u, err := url.Parse(in.URL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return in, ErrInvalidMaterial
		}
	}
	return in, nil
}
