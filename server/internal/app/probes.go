package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/setup"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func probes(cfg config.Config, pool *pgxpool.Pool, ext Integrations) []setup.Probe {
	storage, _ := ext.Storage.(pinger)
	out := []setup.Probe{
		probe{name: "database", required: true, check: pool.Ping},
		probe{name: "storage", required: true, check: when(storage != nil, func(ctx context.Context) error { return storage.Ping(ctx) })},
		probe{name: "mail", check: when(cfg.Mail.Enabled(), func(ctx context.Context) error { return mail.Ping(ctx, cfg.Mail) })},
		probe{name: "push", check: when(ext.Pusher != nil, func(context.Context) error { return nil })},
		probe{name: "model", check: model(cfg.LLM)},
		probe{name: "embeddings", check: model(cfg.Embed)},
		probe{name: "reranker", check: model(cfg.Rerank)},
		probe{name: "speech", check: model(cfg.Speech)},
		probe{name: "calls", check: when(ext.Calls != nil, func(ctx context.Context) error { return ext.Calls.Ping(ctx) })},
		probe{name: "recording", check: when(ext.Calls != nil && ext.Calls.Records(), func(context.Context) error { return nil })},
	}
	return out
}

func model(cfg llm.Config) func(context.Context) error {
	return when(cfg.Enabled(), func(ctx context.Context) error { return llm.Ping(ctx, cfg) })
}

func when(configured bool, check func(context.Context) error) func(context.Context) error {
	if !configured {
		return func(context.Context) error { return setup.ErrNotConfigured }
	}
	return check
}
