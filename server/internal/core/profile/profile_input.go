package profile

import (
	"regexp"
	"unicode/utf8"
)

const (
	maxContacts     = 50
	maxContactValue = 200
)

var contactKindPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,19}$`)

type ProfileInput struct {
	DisplayName string    `json:"displayName"`
	About       string    `json:"about"`
	Contacts    []Contact `json:"contacts"`
}

func (in ProfileInput) validate() error {
	if utf8.RuneCountInString(in.DisplayName) > 100 || utf8.RuneCountInString(in.About) > 2000 || len(in.Contacts) > maxContacts {
		return ErrInvalidProfile
	}
	for _, c := range in.Contacts {
		if !contactKindPattern.MatchString(c.Kind) || c.Value == "" || utf8.RuneCountInString(c.Value) > maxContactValue {
			return ErrInvalidProfile
		}
	}
	return nil
}
