package sandbox

import (
	"strings"
	"unicode/utf8"
)

type NodeInput struct {
	Name    string `json:"name"`
	ClassID string `json:"classId"`
	GPUs    int    `json:"gpus"`
}

func (in NodeInput) normalized() (NodeInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > 100 || in.ClassID == "" || in.GPUs < 1 || in.GPUs > 64 {
		return in, ErrInvalidNode
	}
	return in, nil
}
