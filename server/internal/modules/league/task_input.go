package league

import (
	"strings"
	"time"
	"unicode/utf8"
)

type TaskInput struct {
	Title      string     `json:"title"`
	About      string     `json:"about"`
	Partner    string     `json:"partner"`
	Points     int        `json:"points"`
	MinSize    int        `json:"minSize"`
	MaxSize    int        `json:"maxSize"`
	SelfFormed bool       `json:"selfFormed"`
	DueAt      *time.Time `json:"dueAt"`
	Closed     bool       `json:"closed"`
}

func (in TaskInput) normalized() (TaskInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.About = strings.TrimSpace(in.About)
	in.Partner = strings.TrimSpace(in.Partner)
	if in.Title == "" || utf8.RuneCountInString(in.Title) > 200 || utf8.RuneCountInString(in.About) > 10000 ||
		utf8.RuneCountInString(in.Partner) > 200 || in.Points < 1 || in.Points > maxWeight ||
		in.MinSize < 2 || in.MaxSize < in.MinSize || in.MaxSize > 12 {
		return in, ErrInvalidTask
	}
	return in, nil
}
