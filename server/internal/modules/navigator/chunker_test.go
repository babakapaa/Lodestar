package navigator

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestChunkKeepsPassagesShort(t *testing.T) {
	word := strings.Repeat("я", 2000)
	text := "Короткий абзац.\n\nЕщё один.\n\n" + strings.Repeat("Длинное предложение про клетку. ", 100) + "\n\n" + word
	parts := chunk(text)
	if len(parts) < 3 || !strings.HasPrefix(parts[0], "Короткий абзац.\n\nЕщё один.") {
		t.Fatalf("parts = %d %q", len(parts), parts[0])
	}
	for _, p := range parts[:len(parts)-1] {
		if utf8.RuneCountInString(p) > passageLength {
			t.Fatalf("passage too long: %d", utf8.RuneCountInString(p))
		}
	}
}

func TestSentencesKeepSeparators(t *testing.T) {
	text := "Первое [1]. Второе?\n- пункт\n\nТретье… Конец"
	list := sentences(text)
	if len(list) != 5 || joinSentences(list) != text {
		t.Fatalf("sentences = %+v", list)
	}
	if !grounded("Цитата «Митохондрии вырабатывают энергию» [1].", []Fragment{{Text: "Митохондрии  вырабатывают ЭНЕРГИЮ для клетки"}}) {
		t.Fatal("quote with other spacing and case not found")
	}
}
