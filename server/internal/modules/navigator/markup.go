package navigator

import (
	"regexp"
	"strings"
)

var (
	bareMarkup = regexp.MustCompile(`^\s*(\|[\s|:]*|[-*+]|\d+[.)]|#{1,6}|>)\s*$`)
	tableRule  = regexp.MustCompile(`^\s*\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)*\|?\s*$`)
	listNumber = regexp.MustCompile(`(^|\n)[ \t]*\d+$`)
	headingRow = regexp.MustCompile(`^\s*#{1,6}\s`)
)

func listMarker(runes []rune, i int) bool {
	if runes[i] != '.' && runes[i] != ')' {
		return false
	}
	start := max(0, i-12)
	return listNumber.MatchString(string(runes[start:i]))
}

func keepSentences(all []sentence, drop []bool) []sentence {
	var out []sentence
	for i, s := range all {
		if !drop[i] {
			out = append(out, s)
			continue
		}
		if len(out) > 0 && strings.Count(s.tail, "\n") > strings.Count(out[len(out)-1].tail, "\n") {
			out[len(out)-1].tail = s.tail
		}
	}
	return out
}

func tidyMarkup(text string) string {
	lines := strings.Split(text, "\n")
	var out []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(strings.TrimSpace(line), "|") && i+1 < len(lines) && tableRule.MatchString(lines[i+1]) && strings.Contains(lines[i+1], "|") {
			j := i + 2
			var rows []string
			for j < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[j]), "|") {
				if !bareMarkup.MatchString(lines[j]) {
					rows = append(rows, lines[j])
				}
				j++
			}
			if len(rows) > 0 {
				out = append(out, line, lines[i+1])
				out = append(out, rows...)
			}
			i = j - 1
			continue
		}
		if bareMarkup.MatchString(line) && !tableRule.MatchString(line) {
			continue
		}
		out = append(out, line)
	}
	text = strings.Join(emptySections(out), "\n")
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(text)
}

func emptySections(lines []string) []string {
	var out []string
	for i, line := range lines {
		if headingRow.MatchString(line) && !sectionBody(lines[i+1:]) {
			continue
		}
		out = append(out, line)
	}
	return out
}

func sectionBody(rest []string) bool {
	for _, line := range rest {
		if headingRow.MatchString(line) {
			return false
		}
		if strings.TrimSpace(line) != "" && !tableRule.MatchString(line) {
			return true
		}
	}
	return false
}
