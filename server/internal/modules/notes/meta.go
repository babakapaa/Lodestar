package notes

import (
	"strings"
	"unicode/utf8"
)

const (
	maxPath  = 1000
	maxTags  = 200
	maxLinks = 2000
	maxLabel = 300
)

type Meta struct {
	Path         string   `json:"path"`
	Tags         []string `json:"tags"`
	Links        []string `json:"links"`
	AssignmentID *string  `json:"assignmentId"`
	LessonID     *string  `json:"lessonId"`
}

func (m Meta) normalized() (Meta, error) {
	m.Path = strings.Trim(strings.TrimSpace(m.Path), "/")
	if m.Path == "" || utf8.RuneCountInString(m.Path) > maxPath || strings.Contains(m.Path, "//") ||
		len(m.Tags) > maxTags || len(m.Links) > maxLinks {
		return m, ErrInvalidNote
	}
	for _, list := range [][]string{m.Tags, m.Links} {
		for _, item := range list {
			if item == "" || utf8.RuneCountInString(item) > maxLabel {
				return m, ErrInvalidNote
			}
		}
	}
	if m.Tags == nil {
		m.Tags = []string{}
	}
	if m.Links == nil {
		m.Links = []string{}
	}
	return m, nil
}
