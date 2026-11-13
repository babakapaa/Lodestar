package sandbox

import (
	"strings"
	"unicode/utf8"
)

const (
	defaultSessionHours = 4
	defaultDemoHours    = 2
)

type RunInput struct {
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	Command    string  `json:"command"`
	TemplateID string  `json:"templateId"`
	ClassID    string  `json:"classId"`
	MaxHours   float64 `json:"maxHours"`
	WorkKind   string  `json:"workKind"`
	WorkID     string  `json:"workId"`
	TeamID     string  `json:"teamId"`
}

func (in RunInput) normalized() (RunInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Command = strings.TrimSpace(in.Command)
	in.MaxHours = hours(in.MaxHours)
	in.MaxHours = defaultHours(in.Kind, in.MaxHours)
	switch {
	case in.Kind != RunJupyter && in.Kind != RunVSCode && in.Kind != RunJob && in.Kind != RunDemo,
		in.Title == "" || utf8.RuneCountInString(in.Title) > 200,
		commanded(in.Kind) && in.Command == "",
		interactive(in.Kind) && in.Command != "",
		len(in.Command) > 2000,
		in.TemplateID == "" || in.ClassID == "",
		in.MaxHours <= 0 || in.MaxHours > 168,
		(in.WorkKind == "") != (in.WorkID == ""):
		return in, ErrInvalidRun
	}
	return in, nil
}

func (in RunInput) cost(rate float64) float64 {
	return hours(in.MaxHours * rate)
}

func defaultHours(kind string, h float64) float64 {
	switch {
	case h != 0:
		return h
	case interactive(kind):
		return defaultSessionHours
	case kind == RunDemo:
		return defaultDemoHours
	}
	return 0
}
