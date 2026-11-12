package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	shutdown, err := durationEnv("LODESTAR_SHUTDOWN_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, err
	}
	return Config{
		HTTPAddr:        stringEnv("LODESTAR_HTTP_ADDR", ":8080"),
		ShutdownTimeout: shutdown,
	}, nil
}

func stringEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
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
