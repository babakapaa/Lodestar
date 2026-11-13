package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	PublicURL       string
	SecretKey       string
	AppRedirects    []string
	WebOrigins      []string
	AppURL          string
	Storage         objectstore.Config
	Mail            mail.Config
	FCMCredentials  string
	LLM             llm.Config
	AccessTTL       time.Duration
	RefreshTTL      time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:     stringEnv("LODESTAR_HTTP_ADDR", ":8080"),
		DatabaseURL:  stringEnv("LODESTAR_DATABASE_URL", ""),
		PublicURL:    strings.TrimRight(stringEnv("LODESTAR_PUBLIC_URL", "http://localhost:8080"), "/"),
		SecretKey:    stringEnv("LODESTAR_SECRET_KEY", ""),
		AppRedirects: listEnv("LODESTAR_APP_REDIRECTS", "io.github.babakapaa.lodestar:/auth,http://localhost:43823/auth"),
		WebOrigins:   listEnv("LODESTAR_WEB_ORIGINS", ""),
		Storage: objectstore.Config{
			Endpoint:       stringEnv("LODESTAR_S3_ENDPOINT", ""),
			PublicEndpoint: stringEnv("LODESTAR_S3_PUBLIC_ENDPOINT", ""),
			Region:         stringEnv("LODESTAR_S3_REGION", "us-east-1"),
			Bucket:         stringEnv("LODESTAR_S3_BUCKET", "lodestar"),
			AccessKey:      stringEnv("LODESTAR_S3_ACCESS_KEY", ""),
			SecretKey:      stringEnv("LODESTAR_S3_SECRET_KEY", ""),
		},
		Mail: mail.Config{
			Addr:     stringEnv("LODESTAR_SMTP_ADDR", ""),
			Username: stringEnv("LODESTAR_SMTP_USERNAME", ""),
			Password: stringEnv("LODESTAR_SMTP_PASSWORD", ""),
			From:     stringEnv("LODESTAR_MAIL_FROM", ""),
		},
		FCMCredentials: stringEnv("LODESTAR_FCM_CREDENTIALS_FILE", ""),
		LLM: llm.Config{
			URL:   stringEnv("LODESTAR_LLM_URL", ""),
			Model: stringEnv("LODESTAR_LLM_MODEL", ""),
			Key:   stringEnv("LODESTAR_LLM_KEY", ""),
		},
	}
	cfg.AppURL = strings.TrimRight(stringEnv("LODESTAR_APP_URL", cfg.PublicURL), "/")
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("LODESTAR_DATABASE_URL is required")
	}
	if cfg.SecretKey == "" {
		return Config{}, errors.New("LODESTAR_SECRET_KEY is required")
	}
	var err error
	if cfg.AccessTTL, err = durationEnv("LODESTAR_ACCESS_TTL", 15*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.RefreshTTL, err = durationEnv("LODESTAR_REFRESH_TTL", 60*24*time.Hour); err != nil {
		return Config{}, err
	}
	if cfg.ShutdownTimeout, err = durationEnv("LODESTAR_SHUTDOWN_TIMEOUT", 15*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.LLM.Timeout, err = durationEnv("LODESTAR_LLM_TIMEOUT", 2*time.Minute); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func stringEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func listEnv(key, fallback string) []string {
	var list []string
	for item := range strings.SplitSeq(stringEnv(key, fallback), ",") {
		if item = strings.TrimSpace(item); item != "" {
			list = append(list, item)
		}
	}
	return list
}

func durationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	return parsed, nil
}
