package study

import "strings"

type CourseInput struct {
	UnitID      *string `json:"unitId"`
	Title       string  `json:"title"`
	Code        *string `json:"code"`
	Description string  `json:"description"`
}

func (in CourseInput) normalized() CourseInput {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	if in.Code != nil {
		code := strings.TrimSpace(*in.Code)
		in.Code = &code
		if code == "" {
			in.Code = nil
		}
	}
	if in.UnitID != nil && *in.UnitID == "" {
		in.UnitID = nil
	}
	return in
}
