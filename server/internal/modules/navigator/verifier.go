package navigator

import (
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Verifier struct {
	model llm.Model
}

func (v Verifier) Check(ctx context.Context, draft string, prompt Prompt) Verdict {
	var kept []sentence
	removed := 0
	for _, s := range sentences(draft) {
		if grounded(s.text, prompt.Fragments) {
			kept = append(kept, s)
		} else {
			removed++
		}
	}
	reply := v.ask(ctx, kept, prompt)
	var final []sentence
	for i, s := range kept {
		if slices.Contains(reply.Unsupported, i+1) {
			removed++
			continue
		}
		final = append(final, s)
	}
	text := joinSentences(final)
	leak := reply.Leak && (prompt.Brief == nil || prompt.Brief.Mode != module.ModeProject)
	return Verdict{Text: text, Citations: citations(text, prompt.Fragments), Leak: leak, Removed: removed}
}

func (v Verifier) ask(ctx context.Context, kept []sentence, prompt Prompt) checkReply {
	if len(kept) == 0 {
		return checkReply{}
	}
	var b strings.Builder
	b.WriteString(fragmentsText(prompt.Fragments))
	b.WriteString("\n\n" + briefText(prompt.Brief, prompt.Now))
	b.WriteString("\n\nСтупень лестницы: " + string(prompt.Step))
	b.WriteString("\n\nПредложения ответа:")
	for i, s := range kept {
		b.WriteString("\n" + strconv.Itoa(i+1) + ". " + s.text)
	}
	raw, err := v.model.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: verifierPrompt},
			{Role: llm.RoleUser, Content: b.String()},
		},
		MaxTokens: 200,
		JSON:      true,
	})
	if err != nil {
		slog.Warn("verify navigator answer", "error", err)
		return checkReply{}
	}
	var out checkReply
	if json.Unmarshal([]byte(jsonObject(raw)), &out) != nil {
		slog.Warn("verify navigator answer", "error", "reply is not JSON")
		return checkReply{}
	}
	return out
}

func grounded(text string, fragments []Fragment) bool {
	refs := cited(text)
	for _, n := range refs {
		if n < 1 || n > len(fragments) {
			return false
		}
	}
	for _, q := range quotes(text) {
		if !quoted(plain(q), refs, fragments) {
			return false
		}
	}
	return true
}

func quoted(q string, refs []int, fragments []Fragment) bool {
	for i, f := range fragments {
		if (len(refs) == 0 || slices.Contains(refs, i+1)) && strings.Contains(plain(f.Text), q) {
			return true
		}
	}
	return false
}

func citations(text string, fragments []Fragment) []Citation {
	var out []Citation
	for _, n := range cited(text) {
		if slices.ContainsFunc(out, func(c Citation) bool { return c.N == n }) {
			continue
		}
		f := fragments[n-1]
		out = append(out, Citation{N: n, MaterialID: f.MaterialID, Title: f.Title, Reference: f.Reference, URL: f.URL, Text: f.Text})
	}
	slices.SortFunc(out, func(a, b Citation) int { return a.N - b.N })
	return out
}

func jsonObject(raw string) string {
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}
