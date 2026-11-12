package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_HTTP_ADDR", "")
	t.Setenv("LODESTAR_SHUTDOWN_TIMEOUT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 15*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 15s", cfg.ShutdownTimeout)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_HTTP_ADDR", ":9000")
	t.Setenv("LODESTAR_SHUTDOWN_TIMEOUT", "3s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPAddr != ":9000" {
		t.Errorf("HTTPAddr = %q, want :9000", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 3*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 3s", cfg.ShutdownTimeout)
	}
}

func TestLoadRequiresDatabase(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load accepted a missing database url")
	}
}

func TestLoadRejectsBadDuration(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SHUTDOWN_TIMEOUT", "soon")

	if _, err := Load(); err == nil {
		t.Fatal("Load accepted an invalid duration")
	}
}
