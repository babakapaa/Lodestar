package team

import (
	"strings"
	"time"
	"unicode/utf8"
)

type Card struct {
	Interests string     `json:"interests"`
	Skills    string     `json:"skills"`
	Hours     string     `json:"hours"`
	Role      string     `json:"role"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func (c Card) normalized() (Card, error) {
	c.Interests = strings.TrimSpace(c.Interests)
	c.Skills = strings.TrimSpace(c.Skills)
	c.Hours = strings.TrimSpace(c.Hours)
	c.Role = strings.TrimSpace(c.Role)
	if utf8.RuneCountInString(c.Interests) > 500 || utf8.RuneCountInString(c.Skills) > 500 ||
		utf8.RuneCountInString(c.Hours) > 200 || utf8.RuneCountInString(c.Role) > 100 {
		return c, ErrInvalidTeam
	}
	c.UpdatedAt = nil
	return c, nil
}

func (c Card) empty() bool {
	return c.Interests == "" && c.Skills == "" && c.Hours == "" && c.Role == ""
}
