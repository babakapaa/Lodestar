package team

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Matcher struct {
	model llm.Model
}

type candidate struct {
	person identity.Person
	card   Card
}

func (m Matcher) Match(ctx context.Context, people []candidate, minSize, maxSize int, must string) ([]SuggestedTeam, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "min: %d\nmax: %d\n", minSize, maxSize)
	if must != "" {
		fmt.Fprintf(&b, "must: %s\n", must)
	}
	byID := make(map[string]identity.Person, len(people))
	for _, c := range people {
		byID[c.person.UserID] = c.person
		fmt.Fprintf(&b, "\nid: %s\nинтересы: %s\nнавыки: %s\nвремя: %s\nроль: %s\n", c.person.UserID, c.card.Interests, c.card.Skills, c.card.Hours, c.card.Role)
	}
	raw, err := m.model.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: matcherPrompt},
			{Role: llm.RoleUser, Content: b.String()},
		},
		Temperature: 0.3,
		MaxTokens:   1500,
		JSON:        true,
	})
	if err != nil {
		return nil, ErrAIFailed
	}
	var out struct {
		Teams []struct {
			Members []string `json:"members"`
			Why     string   `json:"why"`
		} `json:"teams"`
	}
	if json.Unmarshal([]byte(objectIn(raw)), &out) != nil {
		return nil, ErrAIFailed
	}
	used := map[string]bool{}
	var teams []SuggestedTeam
	for _, t := range out.Teams {
		var members []identity.Person
		for _, id := range t.Members {
			p, ok := byID[id]
			if !ok || used[id] {
				continue
			}
			used[id] = true
			members = append(members, p)
		}
		if len(members) < minSize || len(members) > maxSize {
			for _, p := range members {
				delete(used, p.UserID)
			}
			continue
		}
		if must != "" && !slices.ContainsFunc(members, func(p identity.Person) bool { return p.UserID == must }) {
			continue
		}
		teams = append(teams, SuggestedTeam{Members: members, Why: strings.TrimSpace(t.Why)})
		if must != "" {
			break
		}
	}
	if len(teams) == 0 {
		return nil, ErrAIFailed
	}
	return teams, nil
}

func objectIn(raw string) string {
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}
