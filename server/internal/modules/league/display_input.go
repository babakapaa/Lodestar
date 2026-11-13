package league

import (
	"strings"
	"unicode/utf8"
)

type DisplayInput struct {
	Display  Display `json:"display"`
	Nickname string  `json:"nickname"`
}

func (in DisplayInput) normalized() (DisplayInput, error) {
	in.Nickname = strings.TrimSpace(in.Nickname)
	if !in.Display.valid() || utf8.RuneCountInString(in.Nickname) > 40 || (in.Display == DisplayNickname && in.Nickname == "") {
		return in, ErrInvalidInput
	}
	return in, nil
}
