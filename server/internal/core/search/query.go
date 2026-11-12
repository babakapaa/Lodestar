package search

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxQueryLength = 200

type Query struct {
	Text  string
	Terms []string
}

func ParseQuery(raw string) (Query, error) {
	raw = strings.TrimSpace(raw)
	if utf8.RuneCountInString(raw) > maxQueryLength {
		return Query{}, ErrInvalidQuery
	}
	terms := strings.FieldsFunc(normalize(raw), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	if len(terms) == 0 {
		return Query{}, ErrInvalidQuery
	}
	return Query{Text: strings.Join(terms, " "), Terms: terms}, nil
}

func normalize(s string) string {
	return strings.NewReplacer("ё", "е", "Ё", "е").Replace(strings.ToLower(s))
}
