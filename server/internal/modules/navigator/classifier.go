package navigator

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Classifier struct {
	model llm.Model
}

func (c Classifier) Classify(ctx context.Context, question string) Classification {
	raw, err := c.model.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: classifierPrompt},
			{Role: llm.RoleUser, Content: question},
		},
		MaxTokens: 50,
		JSON:      true,
	})
	fallback := Classification{Intent: IntentUnderstand}
	if err != nil {
		return fallback
	}
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		raw = raw[start : end+1]
	}
	var out Classification
	if json.Unmarshal([]byte(raw), &out) != nil || !out.Intent.Choosable() {
		return Classification{Intent: IntentUnderstand, Distress: out.Distress}
	}
	return out
}
