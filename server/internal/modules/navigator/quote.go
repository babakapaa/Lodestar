package navigator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var citeRe = regexp.MustCompile(`\[(\d{1,2})\]`)

var quoteRe = regexp.MustCompile(`«([^»]{8,})»`)

func cited(text string) []int {
	var out []int
	for _, m := range citeRe.FindAllStringSubmatch(text, -1) {
		n, _ := strconv.Atoi(m[1])
		out = append(out, n)
	}
	return out
}

func quotes(text string) []string {
	var out []string
	for _, m := range quoteRe.FindAllStringSubmatch(text, -1) {
		out = append(out, m[1])
	}
	return out
}

func plain(text string) string {
	var b strings.Builder
	space := false
	for _, r := range strings.ToLower(text) {
		switch {
		case r == 'ё':
			r = 'е'
		case strings.ContainsRune("‐‑‒–—―", r):
			r = '-'
		case strings.ContainsRune("“”„\"'’", r):
			continue
		}
		if unicode.IsSpace(r) {
			space = true
			continue
		}
		if space && b.Len() > 0 {
			b.WriteByte(' ')
		}
		space = false
		b.WriteRune(r)
	}
	return b.String()
}
