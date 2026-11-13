package navigator

import "strings"

type DialogInput struct {
	ID        string `json:"id"`
	ScopeKind string `json:"scopeKind"`
	ScopeID   string `json:"scopeId"`
}

func (in DialogInput) normalized() (DialogInput, error) {
	in.ScopeKind = strings.TrimSpace(in.ScopeKind)
	in.ScopeID = strings.TrimSpace(in.ScopeID)
	if !validID(in.ID) || len(in.ScopeKind) > 30 || len(in.ScopeID) > 100 || (in.ScopeKind == "") != (in.ScopeID == "") {
		return in, ErrInvalidDialog
	}
	return in, nil
}
