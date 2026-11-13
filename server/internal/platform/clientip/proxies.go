package clientip

import (
	"fmt"
	"net/http"
	"net/netip"
	"strings"
)

type Proxies []netip.Prefix

func Parse(list []string) (Proxies, error) {
	var out Proxies
	for _, item := range list {
		prefix, err := ParsePrefix(item)
		if err != nil {
			return nil, err
		}
		out = append(out, prefix)
	}
	return out, nil
}

func ParsePrefix(value string) (netip.Prefix, error) {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "/") {
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			return netip.Prefix{}, fmt.Errorf("network %q: %w", value, err)
		}
		return prefix.Masked(), nil
	}
	addr, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("network %q: %w", value, err)
	}
	return netip.PrefixFrom(addr.Unmap(), addr.Unmap().BitLen()), nil
}

func Of(r *http.Request) (netip.Addr, bool) {
	ap, err := netip.ParseAddrPort(r.RemoteAddr)
	if err != nil {
		return netip.Addr{}, false
	}
	return ap.Addr().Unmap(), true
}

func (p Proxies) Wrap(next http.Handler) http.Handler {
	if len(p) == 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if addr, ok := p.resolve(r); ok {
			r = r.Clone(r.Context())
			r.RemoteAddr = netip.AddrPortFrom(addr, 0).String()
		}
		next.ServeHTTP(w, r)
	})
}

func (p Proxies) resolve(r *http.Request) (netip.Addr, bool) {
	peer, ok := Of(r)
	if !ok || !p.trusted(peer) {
		return netip.Addr{}, false
	}
	hops := strings.Split(strings.Join(r.Header.Values("X-Forwarded-For"), ","), ",")
	for i := len(hops) - 1; i >= 0; i-- {
		addr, err := netip.ParseAddr(strings.TrimSpace(hops[i]))
		if err != nil {
			return netip.Addr{}, false
		}
		if addr = addr.Unmap(); !p.trusted(addr) {
			return addr, true
		}
	}
	return netip.Addr{}, false
}

func (p Proxies) trusted(addr netip.Addr) bool {
	for _, prefix := range p {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
