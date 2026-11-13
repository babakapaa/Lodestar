package study

import (
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type CourseInput struct {
	UnitID        *string `json:"unitId"`
	Title         string  `json:"title"`
	Code          *string `json:"code"`
	Description   string  `json:"description"`
	NavigatorMode string  `json:"navigatorMode"`
}

func (in CourseInput) normalized() CourseInput {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	if in.NavigatorMode == "" {
		in.NavigatorMode = module.ModeStudy
	}
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
