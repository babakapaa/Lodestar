package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SECRET_KEY", "key")
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
	if cfg.AccessTTL != 15*time.Minute {
		t.Errorf("AccessTTL = %v, want 15m", cfg.AccessTTL)
	}
	if len(cfg.AppRedirects) != 2 {
		t.Errorf("AppRedirects = %v, want mobile and desktop defaults", cfg.AppRedirects)
	}
	if cfg.Tester || cfg.WebDir != "" {
		t.Errorf("Tester = %v, WebDir = %q, want off and empty", cfg.Tester, cfg.WebDir)
	}
}

func TestLoadTesterMode(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SECRET_KEY", "key")
	t.Setenv("LODESTAR_TESTER_MODE", "On")
	t.Setenv("LODESTAR_WEB_DIR", "/web")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Tester || cfg.WebDir != "/web" {
		t.Errorf("Tester = %v, WebDir = %q", cfg.Tester, cfg.WebDir)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SECRET_KEY", "key")
	t.Setenv("LODESTAR_HTTP_ADDR", ":9000")
	t.Setenv("LODESTAR_SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("LODESTAR_PUBLIC_URL", "https://lodestar.example/")
	t.Setenv("LODESTAR_APP_REDIRECTS", "app:/auth, http://localhost:5000/auth")

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
	if cfg.PublicURL != "https://lodestar.example" {
		t.Errorf("PublicURL = %q", cfg.PublicURL)
	}
	if len(cfg.AppRedirects) != 2 || cfg.AppRedirects[1] != "http://localhost:5000/auth" {
		t.Errorf("AppRedirects = %v", cfg.AppRedirects)
	}
}

func TestLoadRequiresDatabase(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load accepted a missing database url")
	}
}

func TestLoadRequiresSecretKey(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SECRET_KEY", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load accepted a missing secret key")
	}
}

func TestLoadRejectsBadDuration(t *testing.T) {
	t.Setenv("LODESTAR_DATABASE_URL", "postgres://db")
	t.Setenv("LODESTAR_SECRET_KEY", "key")
	t.Setenv("LODESTAR_SHUTDOWN_TIMEOUT", "soon")

	if _, err := Load(); err == nil {
		t.Fatal("Load accepted an invalid duration")
	}
}
