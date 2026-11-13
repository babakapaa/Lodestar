package study

import (
	"fmt"
	"strings"
)

const (
	maxMaterialWork = 12000
	maxMaterialText = 3000
)

type DefenseMaterial struct {
	Title       string
	Description string
	Work        string
	Files       []string
	Map         *SearchMap
}

func (m DefenseMaterial) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Задание: %s\n", m.Title)
	if m.Description != "" {
		fmt.Fprintf(&b, "Описание задания:\n%s\n", clip(m.Description, maxMaterialText))
	}
	if m.Work != "" {
		fmt.Fprintf(&b, "\nТекст работы:\n%s\n", clip(m.Work, maxMaterialWork))
	}
	if len(m.Files) > 0 {
		fmt.Fprintf(&b, "\nПриложенные файлы: %s\n", strings.Join(m.Files, ", "))
	}
	if m.Map != nil {
		b.WriteString("\nКарта поиска.\n")
		for _, list := range []struct {
			name    string
			sources []MapSource
		}{{"Использованы", m.Map.Used}, {"Отброшены", m.Map.Rejected}} {
			for _, s := range list.sources {
				fmt.Fprintf(&b, "- %s: %s. Где найден: %s. Почему: %s\n", list.name, s.Reference, s.Found, s.Why)
			}
		}
		if len(m.Map.Queries) > 0 {
			fmt.Fprintf(&b, "Запросы: %s\n", strings.Join(m.Map.Queries, "; "))
		}
		if m.Map.Path != "" {
			fmt.Fprintf(&b, "Путь поиска:\n%s\n", clip(m.Map.Path, maxMaterialText))
		}
	}
	return b.String()
}

func clip(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit]) + "…"
}
