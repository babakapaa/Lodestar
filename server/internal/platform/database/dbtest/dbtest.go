package dbtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const urlEnv = "LODESTAR_TEST_DATABASE_URL"

func New(t *testing.T) *pgxpool.Pool {
	t.Helper()
	adminURL := os.Getenv(urlEnv)
	if adminURL == "" {
		if os.Getenv("CI") != "" {
			t.Fatalf("%s must be set in CI", urlEnv)
		}
		t.Skipf("%s is not set", urlEnv)
	}
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, adminURL)
	if err != nil {
		t.Fatalf("connect admin: %v", err)
	}
	t.Cleanup(func() { _ = admin.Close(context.Background()) })

	name := "lodestar_test_" + randomSuffix(t)
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatalf("create database: %v", err)
	}
	testURL, err := url.Parse(adminURL)
	if err != nil {
		t.Fatalf("%s must be a postgres:// URL: %v", urlEnv, err)
	}
	testURL.Path = "/" + name
	pool, err := database.Open(ctx, testURL.String())
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
		if _, err := admin.Exec(context.Background(), "DROP DATABASE "+name+" WITH (FORCE)"); err != nil {
			t.Errorf("drop database: %v", err)
		}
	})
	return pool
}

func randomSuffix(t *testing.T) string {
	t.Helper()
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("random: %v", err)
	}
	return hex.EncodeToString(b)
}
