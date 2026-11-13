package sandbox

import (
	"strings"
	"unicode/utf8"
)

type TemplateInput struct {
	Name     string `json:"name"`
	About    string `json:"about"`
	Image    string `json:"image"`
	Archived bool   `json:"archived"`
}

func (in TemplateInput) normalized() (TemplateInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.About = strings.TrimSpace(in.About)
	in.Image = strings.TrimSpace(in.Image)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > 100 || utf8.RuneCountInString(in.About) > 1000 ||
		in.Image == "" || len(in.Image) > 300 || strings.ContainsAny(in.Image, " \t\n") {
		return in, ErrInvalidTemplate
	}
	return in, nil
}
