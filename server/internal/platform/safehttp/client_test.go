package safehttp_test

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/safehttp"
)

func TestPrivateAddressesAreRefused(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	resp, err := safehttp.NewClient(time.Second).Get(srv.URL)
	if err == nil {
		_ = resp.Body.Close()
	}
	if !errors.Is(err, safehttp.ErrPrivateAddress) {
		t.Fatalf("loopback was reachable: %v", err)
	}
}

func TestPublicRanges(t *testing.T) {
	for addr, public := range map[string]bool{
		"8.8.8.8": true, "2a00:1450::1": true, "10.1.2.3": false, "192.168.0.1": false,
		"169.254.169.254": false, "100.64.0.1": false, "::1": false, "fc00::1": false, "0.0.0.0": false,
	} {
		if got := safehttp.Public(net.ParseIP(addr)); got != public {
			t.Errorf("%s: public=%v", addr, got)
		}
	}
}
