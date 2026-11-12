package safehttp

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

var ErrPrivateAddress = errors.New("address is not public")

var cgnat = net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}

func NewClient(timeout time.Duration) *http.Client {
	dialer := &net.Dialer{Timeout: 5 * time.Second, Control: control}
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext:            dialer.DialContext,
			TLSHandshakeTimeout:    5 * time.Second,
			ResponseHeaderTimeout:  timeout,
			MaxResponseHeaderBytes: 1 << 16,
		},
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

func Public(ip net.IP) bool {
	blocked := ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsMulticast() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() ||
		cgnat.Contains(ip)
	return !blocked
}

func control(_, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("dial %s: %w", address, err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !Public(ip) {
		return fmt.Errorf("dial %s: %w", address, ErrPrivateAddress)
	}
	return nil
}
