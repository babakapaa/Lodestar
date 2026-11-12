package notes

import "unicode/utf8"

const maxUpdateBytes = 8 << 20

type UpdateInput struct {
	Key  string `json:"key"`
	Data []byte `json:"data"`
}

func (in UpdateInput) valid() bool {
	n := utf8.RuneCountInString(in.Key)
	return n >= 1 && n <= 100 && len(in.Data) > 0 && len(in.Data) <= maxUpdateBytes
}
