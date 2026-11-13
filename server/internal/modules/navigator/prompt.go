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
	Fragments []string
	Memory    []Message
	Question  string
	Now       time.Time
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
	if len(p.Fragments) == 0 {
		parts = append(parts, "Фрагменты: по этому вопросу в базе ничего не найдено.")
	} else {
		var fragments strings.Builder
		fragments.WriteString("Фрагменты:")
		for i, f := range p.Fragments {
			fragments.WriteString("\n[" + strconv.Itoa(i+1) + "] " + f)
		}
		parts = append(parts, fragments.String())
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
