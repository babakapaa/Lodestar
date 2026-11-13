package clientip_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/clientip"
)

func TestProxiesTrustOnlyTheirOwnHops(t *testing.T) {
	proxies, err := clientip.Parse([]string{"10.0.0.0/8", "192.168.1.5"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := clientip.Parse([]string{"nonsense"}); err == nil {
		t.Fatal("bad network accepted")
	}
	cases := []struct {
		peer, forwarded, want string
	}{
		{"10.1.2.3:4000", "203.0.113.7, 198.51.100.2, 10.0.0.9", "198.51.100.2"},
		{"192.168.1.5:80", "203.0.113.7", "203.0.113.7"},
		{"203.0.113.50:5000", "1.2.3.4", "203.0.113.50"},
		{"10.1.2.3:4000", "", "10.1.2.3"},
		{"10.1.2.3:4000", "garbage", "10.1.2.3"},
	}
	for _, c := range cases {
		var got string
		h := proxies.Wrap(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			addr, _ := clientip.Of(r)
			got = addr.String()
		}))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = c.peer
		if c.forwarded != "" {
			req.Header.Set("X-Forwarded-For", c.forwarded)
		}
		h.ServeHTTP(httptest.NewRecorder(), req)
		if got != c.want {
			t.Errorf("%s via %q = %s, want %s", c.peer, c.forwarded, got, c.want)
		}
	}
}
