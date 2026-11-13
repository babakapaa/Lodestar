package team

import (
	"strings"
	"unicode/utf8"
)

const maxName = 100

type TeamInput struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func (in TeamInput) normalized() (TeamInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > maxName {
		return in, ErrInvalidTeam
	}
	return in, nil
}
