package study

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

const (
	minDraftQuestions = 5
	maxDraftQuestions = 7
)

type DefenseDrafter struct {
	model llm.Model
}

func (d DefenseDrafter) ready() bool {
	return d.model != nil
}

func (d DefenseDrafter) Questions(ctx context.Context, m DefenseMaterial) ([]DefenseQuestion, error) {
	raw, err := d.model.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: questionsPrompt},
			{Role: llm.RoleUser, Content: m.Text()},
		},
		Temperature: 0.4,
		MaxTokens:   1200,
		JSON:        true,
	})
	if err != nil {
		return nil, ErrAIFailed
	}
	var out struct {
		Questions []DefenseQuestion `json:"questions"`
	}
	if json.Unmarshal([]byte(objectIn(raw)), &out) != nil {
		return nil, ErrAIFailed
	}
	questions := make([]DefenseQuestion, 0, maxDraftQuestions)
	for _, q := range out.Questions {
		q.Mark, q.Note = nil, ""
		if q, ok := q.normalized(); ok && len(questions) < maxDraftQuestions {
			questions = append(questions, q)
		}
	}
	if len(questions) == 0 {
		return nil, ErrAIFailed
	}
	return questions, nil
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

func objectIn(raw string) string {
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}
