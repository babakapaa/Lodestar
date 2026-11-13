package league

import (
	"strings"
	"unicode/utf8"
)

type AwardInput struct {
	Title    string `json:"title"`
	About    string `json:"about"`
	GPUHours int    `json:"gpuHours"`
	Quantity *int   `json:"quantity"`
	Archived bool   `json:"archived"`
}

func (in AwardInput) normalized() (AwardInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.About = strings.TrimSpace(in.About)
	if in.Title == "" || utf8.RuneCountInString(in.Title) > 200 || utf8.RuneCountInString(in.About) > 2000 ||
		in.GPUHours < 0 || in.GPUHours > 100000 || (in.Quantity != nil && (*in.Quantity < 1 || *in.Quantity > 100000)) {
		return in, ErrInvalidAward
	}
	return in, nil
}
