package navigator

import (
	"strings"
	"unicode/utf8"
)

const passageLength = 1200

func chunk(text string) []string {
	var out []string
	var current strings.Builder
	flush := func() {
		if s := strings.TrimSpace(current.String()); s != "" {
			out = append(out, s)
		}
		current.Reset()
	}
	for _, paragraph := range strings.Split(text, "\n\n") {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		for _, piece := range split(paragraph) {
			if current.Len() > 0 && utf8.RuneCountInString(current.String())+utf8.RuneCountInString(piece)+2 > passageLength {
				flush()
			}
			if current.Len() > 0 {
				current.WriteString("\n\n")
			}
			current.WriteString(piece)
		}
	}
	flush()
	return out
}

func split(paragraph string) []string {
	if utf8.RuneCountInString(paragraph) <= passageLength {
		return []string{paragraph}
	}
	var out []string
	var current strings.Builder
	for _, sentence := range sentences(paragraph) {
		for _, word := range wrap(sentence.text) {
			if current.Len() > 0 && utf8.RuneCountInString(current.String())+utf8.RuneCountInString(word)+1 > passageLength {
				out = append(out, strings.TrimSpace(current.String()))
				current.Reset()
			}
			if current.Len() > 0 {
				current.WriteString(" ")
			}
			current.WriteString(word)
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		out = append(out, s)
	}
	return out
}

func wrap(sentence string) []string {
	if utf8.RuneCountInString(sentence) <= passageLength {
		return []string{sentence}
	}
	return strings.Fields(sentence)
}
