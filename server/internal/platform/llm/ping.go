package llm

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func Ping(ctx context.Context, cfg Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(cfg.URL, "/")+"/models", nil)
	if err != nil {
		return err
	}
	if cfg.Key != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Key)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrModelFailed, err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden || resp.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("%w: status %d", ErrModelFailed, resp.StatusCode)
	}
	return nil
}
