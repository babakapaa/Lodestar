package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	PublicURL       string
	SecretKey       string
	AppRedirects    []string
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
		AppRedirects: listEnv("LODESTAR_APP_REDIRECTS", "io.github.babakapaa.lodestar:/auth"),
	}
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
