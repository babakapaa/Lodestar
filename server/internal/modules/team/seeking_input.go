package team

import (
	"strings"
	"unicode/utf8"
)

const maxNote = 500

type SeekingInput struct {
	Seeking bool   `json:"seeking"`
	Note    string `json:"note"`
}

func (in SeekingInput) normalized() (SeekingInput, error) {
	in.Note = strings.TrimSpace(in.Note)
	if utf8.RuneCountInString(in.Note) > maxNote {
		return in, ErrInvalidTeam
	}
	return in, nil
}
