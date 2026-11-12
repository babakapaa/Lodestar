package secretbox

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

func testKey() string {
	return base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, KeySize))
}

func TestSealOpen(t *testing.T) {
	box, err := New(testKey())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	sealed, err := box.Seal([]byte("client-secret"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	plain, err := box.Open(sealed)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if string(plain) != "client-secret" {
		t.Fatalf("plain = %q", plain)
	}
}

func TestOpenRejectsTampering(t *testing.T) {
	box, err := New(testKey())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	sealed, err := box.Seal([]byte("x"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	sealed[len(sealed)-1] ^= 1

	if _, err := box.Open(sealed); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("err = %v, want ErrCorrupt", err)
	}
}

func TestNewRejectsShortKey(t *testing.T) {
	if _, err := New(base64.StdEncoding.EncodeToString([]byte("short"))); !errors.Is(err, ErrBadKey) {
		t.Fatalf("err = %v, want ErrBadKey", err)
	}
}
