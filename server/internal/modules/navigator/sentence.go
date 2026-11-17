package navigator

import (
	"strings"
	"unicode"
)

type sentence struct {
	text string
	tail string
}

func sentences(text string) []sentence {
	var out []sentence
	runes := []rune(text)
	start := 0
	for i := 0; i < len(runes); i++ {
		end := strings.ContainsRune(".!?…", runes[i]) || runes[i] == '\n'
		if !end || (runes[i] != '\n' && i+1 < len(runes) && !unicode.IsSpace(runes[i+1])) || listMarker(runes, i) {
			continue
		}
		j := i + 1
		for j < len(runes) && unicode.IsSpace(runes[j]) {
			j++
		}
		tail := string(runes[i+1 : j])
		if runes[i] == '\n' {
			tail = "\n" + tail
		}
		if body := strings.TrimSpace(string(runes[start : i+1])); body != "" {
			out = append(out, sentence{text: body, tail: tail})
		} else if len(out) > 0 {
			out[len(out)-1].tail += string(runes[start:j])
		}
		start, i = j, j-1
	}
	if body := strings.TrimSpace(string(runes[start:])); body != "" {
		out = append(out, sentence{text: body})
	}
	return out
}

func joinSentences(list []sentence) string {
	var b strings.Builder
	for _, s := range list {
		b.WriteString(s.text)
		b.WriteString(s.tail)
	}
	return strings.TrimSpace(b.String())
}
