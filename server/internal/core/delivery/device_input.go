package delivery

import (
	"slices"
	"strings"
	"unicode/utf8"
)

var platforms = []string{"android", "ios", "web"}

type DeviceInput struct {
	Platform string `json:"platform"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

func (in DeviceInput) normalized() (DeviceInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Language = baseLanguage(in.Language)
	if !slices.Contains(platforms, in.Platform) || in.Token == "" || len(in.Token) > 4096 || utf8.RuneCountInString(in.Name) > 100 || len(in.Language) > 10 {
		return DeviceInput{}, ErrInvalidDevice
	}
	return in, nil
}

func baseLanguage(tag string) string {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if i := strings.IndexAny(tag, "-_"); i >= 0 {
		tag = tag[:i]
	}
	return tag
}
