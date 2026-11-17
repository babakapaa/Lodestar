package navigator

import (
	"strings"
	"testing"
)

func TestPrunedMarkupKeepsStructure(t *testing.T) {
	draft := "### Понятия\n\n| Понятие | Что это |\n|---|---|\n| Объём | Около 15 мл. |\n| Масса | Около 20 г. |\n\n### Как измеряют\n1. Орхидометр Прадера.\n2. УЗИ.\n3. Свободный факт.\n\n### Вопрос\nЧто вы уже знаете?"
	all := sentences(draft)
	drop := make([]bool, len(all))
	for i, s := range all {
		drop[i] = strings.Contains(s.text, "Около") || strings.Contains(s.text, "Свободный")
	}
	got := tidyMarkup(joinSentences(keepSentences(all, drop)))
	want := "### Как измеряют\n1. Орхидометр Прадера.\n2. УЗИ.\n\n### Вопрос\nЧто вы уже знаете?"
	if got != want {
		t.Fatalf("got %q", got)
	}
}
