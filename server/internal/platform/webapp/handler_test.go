package webapp_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/webapp"
)

func TestServesFilesAndIndexForAnyPage(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{"index.html": "<app>", "main.dart.js": "code"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	api := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	h := webapp.New(root, api)

	cases := []struct {
		method, path string
		status       int
		body         string
	}{
		{http.MethodGet, "/", http.StatusOK, "<app>"},
		{http.MethodGet, "/courses/42/assignments", http.StatusOK, "<app>"},
		{http.MethodGet, "/main.dart.js", http.StatusOK, "code"},
		{http.MethodGet, "/v1/meta", http.StatusTeapot, ""},
		{http.MethodGet, "/healthz", http.StatusTeapot, ""},
		{http.MethodPost, "/courses", http.StatusTeapot, ""},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
		if rec.Code != c.status || !strings.Contains(rec.Body.String(), c.body) {
			t.Errorf("%s %s = %d %q, want %d %q", c.method, c.path, rec.Code, rec.Body.String(), c.status, c.body)
		}
	}
}
