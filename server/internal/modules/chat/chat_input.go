package chat

import (
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	maxTitle   = 200
	maxAbout   = 2000
	maxInvited = 200
)

type ChatInput struct {
	Kind      string   `json:"kind"`
	Title     string   `json:"title"`
	About     string   `json:"about"`
	MemberIDs []string `json:"memberIds"`
	Origin    *Origin  `json:"origin"`
}

func (in ChatInput) normalized() (ChatInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.About = strings.TrimSpace(in.About)
	if in.Kind != KindGroup && in.Kind != KindChannel {
		return in, ErrInvalidChat
	}
	if utf8.RuneCountInString(in.Title) > maxTitle || utf8.RuneCountInString(in.About) > maxAbout {
		return in, ErrInvalidChat
	}
	if in.Origin == nil && in.Title == "" {
		return in, ErrInvalidChat
	}
	if in.Origin != nil && (in.Origin.Kind == "" || len(in.Origin.Kind) > 60 || in.Origin.ID == "" || len(in.Origin.ID) > 100) {
		return in, ErrInvalidChat
	}
	ids, err := memberIDs(in.MemberIDs)
	in.MemberIDs = ids
	return in, err
}

func memberIDs(ids []string) ([]string, error) {
	if len(ids) > maxInvited {
		return nil, ErrInvalidChat
	}
	for _, id := range ids {
		if !validID(id) {
			return nil, ErrInvalidChat
		}
	}
	out := slices.Clone(ids)
	slices.Sort(out)
	return slices.Compact(out), nil
}
