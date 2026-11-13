package attendance

import "strings"

type CheckInput struct {
	SessionID string `json:"sessionId"`
	Token     string `json:"token"`
	Code      string `json:"code"`
}

func (in CheckInput) normalized() (CheckInput, error) {
	in.SessionID, in.Token, in.Code = strings.TrimSpace(in.SessionID), strings.TrimSpace(in.Token), strings.TrimSpace(in.Code)
	byQR := in.SessionID != "" && in.Token != ""
	byCode := len(in.Code) == 4 && strings.Trim(in.Code, "0123456789") == ""
	if byQR == byCode || (in.Code != "" && !byCode) {
		return in, ErrInvalidCheck
	}
	return in, nil
}

func (in CheckInput) method() MarkMethod {
	if in.Code != "" {
		return ByCode
	}
	return ByQR
}
