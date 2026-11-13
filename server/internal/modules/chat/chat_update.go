package chat

import (
	"strings"
	"unicode/utf8"
)

type ChatUpdate struct {
	Title string `json:"title"`
	About string `json:"about"`
}

func (in ChatUpdate) normalized() (ChatUpdate, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.About = strings.TrimSpace(in.About)
	if in.Title == "" || utf8.RuneCountInString(in.Title) > maxTitle || utf8.RuneCountInString(in.About) > maxAbout {
		return in, ErrInvalidChat
	}
	return in, nil
}
