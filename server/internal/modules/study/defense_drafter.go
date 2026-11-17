package study

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type DefenseDrafter struct {
	model llm.Model
}

func (d DefenseDrafter) ready() bool {
	return d.model != nil
}

func (d DefenseDrafter) Protocol(ctx context.Context, title string, questions []DefenseQuestion, notes string) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "Задание: %s\n", title)
	for i, q := range questions {
		mark := "нет отметки"
		if q.Mark != nil {
			mark = strconv.Itoa(*q.Mark)
		}
		fmt.Fprintf(&b, "\nВопрос %d (%s): %s\nОтметка: %s\n", i+1, q.Kind, q.Text, mark)
		if q.Note != "" {
			fmt.Fprintf(&b, "Заметка: %s\n", q.Note)
		}
	}
	if notes != "" {
		fmt.Fprintf(&b, "\nОбщие заметки: %s\n", notes)
	}
	raw, err := d.model.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: protocolPrompt},
			{Role: llm.RoleUser, Content: b.String()},
		},
		Temperature: 0.2,
		MaxTokens:   600,
	})
	if err != nil || strings.TrimSpace(raw) == "" {
		return "", ErrAIFailed
	}
	return strings.TrimSpace(raw), nil
}
