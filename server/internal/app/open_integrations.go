package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

func OpenIntegrations(ctx context.Context, cfg config.Config) (Integrations, error) {
	var ext Integrations
	storage, err := OpenStorage(ctx, cfg.Storage)
	if err != nil {
		return Integrations{}, err
	}
	ext.Storage = storage
	if cfg.Mail.Enabled() {
		mailer, err := mail.New(cfg.Mail)
		if err != nil {
			return Integrations{}, err
		}
		ext.Mailer = mailer
	} else {
		slog.Warn("email is not configured, set LODESTAR_SMTP_ADDR and LODESTAR_MAIL_FROM to send notifications by email")
	}
	if cfg.FCMCredentials != "" {
		raw, err := os.ReadFile(cfg.FCMCredentials)
		if err != nil {
			return Integrations{}, fmt.Errorf("LODESTAR_FCM_CREDENTIALS_FILE: %w", err)
		}
		pusher, err := fcm.New(raw)
		if err != nil {
			return Integrations{}, err
		}
		ext.Pusher = pusher
	} else {
		slog.Warn("push is not configured, set LODESTAR_FCM_CREDENTIALS_FILE to send push notifications")
	}
	if cfg.LLM.Enabled() {
		ext.Model = llm.NewClient(cfg.LLM)
	} else {
		slog.Warn("navigator is not configured, set LODESTAR_LLM_URL and LODESTAR_LLM_MODEL to an OpenAI-compatible server such as vLLM or Ollama")
	}
	return ext, nil
}
