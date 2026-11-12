package files

import (
	"mime"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FileInput struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func (in FileInput) normalized() (FileInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || utf8.RuneCountInString(in.Name) > 255 || strings.ContainsFunc(in.Name, unicode.IsControl) || strings.ContainsAny(in.Name, `/\`) {
		return FileInput{}, ErrInvalidFile
	}
	if in.ContentType == "" {
		in.ContentType = "application/octet-stream"
	}
	if _, _, err := mime.ParseMediaType(in.ContentType); err != nil || len(in.ContentType) > 255 {
		return FileInput{}, ErrInvalidFile
	}
	if in.Size <= 0 {
		return FileInput{}, ErrInvalidFile
	}
	return in, nil
}
