package navigator

import (
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func briefText(b *module.Brief, now time.Time) string {
	if b == nil {
		return "Вопрос задан не из задания."
	}
	var out strings.Builder
	if b.Kind == scopeCourse {
		out.WriteString("Курс: " + b.Title)
	} else {
		out.WriteString("Задание: " + b.Title)
		if b.Course != "" {
			out.WriteString("\nКурс: " + b.Course)
		}
	}
	if b.DueAt != nil {
		out.WriteString("\nСрок: " + b.DueAt.UTC().Format(time.RFC3339))
		out.WriteString("\nСейчас: " + now.UTC().Format(time.RFC3339))
	}
	if b.Description != "" {
		out.WriteString("\nОписание:\n" + b.Description)
	}
	return out.String()
}
