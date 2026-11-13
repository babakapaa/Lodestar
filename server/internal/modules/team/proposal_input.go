package team

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

type ProposalInput struct {
	Work json.RawMessage `json:"work"`
	Note string          `json:"note"`
}

func (in ProposalInput) normalized() (ProposalInput, error) {
	in.Note = strings.TrimSpace(in.Note)
	if len(in.Work) == 0 || utf8.RuneCountInString(in.Note) > maxReason {
		return in, ErrInvalidTeam
	}
	return in, nil
}
