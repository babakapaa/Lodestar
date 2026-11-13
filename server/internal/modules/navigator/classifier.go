package navigator

import (
	"context"
	"encoding/json"

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
	var out Classification
	if json.Unmarshal([]byte(jsonObject(raw)), &out) != nil || !out.Intent.Choosable() {
		return Classification{Intent: IntentUnderstand, Distress: out.Distress}
	}
	return out
}
