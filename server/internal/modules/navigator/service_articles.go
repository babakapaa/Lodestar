package navigator

import (
	"context"
	"log/slog"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

func (s Service) articles(ctx context.Context, prompt Prompt) []Article {
	if s.library.Client == nil || !prompt.Library || len(prompt.Fragments) > 0 || !prompt.Step.cites() {
		return nil
	}
	list, err := s.library.Search(ctx, s.libraryQuery(ctx, prompt.query()))
	if err != nil {
		slog.Warn("search open library", "error", err)
		return nil
	}
	return list
}

func (s Service) libraryQuery(ctx context.Context, question string) string {
	if s.k.Model == nil || question == "" {
		return question
	}
	reply, err := s.k.Model.Complete(ctx, llm.Request{Messages: []llm.Message{
		{Role: llm.RoleSystem, Content: libraryPrompt},
		{Role: llm.RoleUser, Content: question},
	}})
	if err != nil {
		return question
	}
	line, _, _ := strings.Cut(strings.TrimSpace(reply), "\n")
	line = strings.Trim(strings.TrimSpace(line), `"«»`)
	if line == "" {
		return question
	}
	return line
}
