package tester

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type localTransport struct {
	prefix  string
	handler http.Handler
	next    http.RoundTripper
}

func (t localTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if !strings.HasPrefix(r.URL.String(), t.prefix) {
		return t.next.RoundTrip(r)
	}
	rec := httptest.NewRecorder()
	t.handler.ServeHTTP(rec, r)
	resp := rec.Result()
	resp.Request = r
	return resp, nil
}
