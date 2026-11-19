package navigator

import (
	"strconv"
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Prompt struct {
	Brief     *module.Brief
	Intent    Intent
	Step      Step
	Fragments []Fragment
	Memory    []Message
	Question  string
	Retry     string
	Now       time.Time
	Library   bool
	Match     int
}

func (p Prompt) Messages() []llm.Message {
	parts := []string{rules}
	if p.Brief != nil && p.Brief.Mode == module.ModeProject {
		parts = append(parts, projectMode)
	} else {
		parts = append(parts, studyMode)
	}
	if text, ok := intents[p.Intent]; ok {
		parts = append(parts, text)
	}
	if text, ok := steps[p.Step]; ok {
		parts = append(parts, text)
	}
	parts = append(parts, briefText(p.Brief, p.Now))
	parts = append(parts, fragmentsText(p.Fragments))
	if p.Retry != "" {
		parts = append(parts, "Прошлый вариант ответа не прошёл проверку: "+p.Retry+" Напиши ответ заново по правилам.")
	}
	out := []llm.Message{{Role: llm.RoleSystem, Content: strings.Join(parts, "\n\n")}}
	for _, m := range p.Memory {
		role := llm.RoleUser
		if m.Role == RoleNavigator {
			role = llm.RoleAssistant
		}
		if m.Text != "" {
			out = append(out, llm.Message{Role: role, Content: m.Text})
		}
	}
	question := p.Question
	if question == "" {
		question = "Дай следующую подсказку."
	}
	return append(out, llm.Message{Role: llm.RoleUser, Content: question})
}

func fragmentsText(fragments []Fragment) string {
	if len(fragments) == 0 {
		return "Фрагменты: по этому вопросу в базе ничего не найдено."
	}
	var b strings.Builder
	b.WriteString("Фрагменты:")
	for i, f := range fragments {
		b.WriteString("\n\n[" + strconv.Itoa(i+1) + "] " + f.Title)
		if f.Reference != "" {
			b.WriteString(" (" + f.Reference + ")")
		}
		b.WriteString("\n" + f.Text)
	}
	return b.String()
}

func (p Prompt) query() string {
	if p.Question != "" {
		return p.Question
	}
	for i := len(p.Memory) - 1; i >= 0; i-- {
		if p.Memory[i].Role == RoleStudent && p.Memory[i].Text != "" {
			return p.Memory[i].Text
		}
	}
	return ""
}
