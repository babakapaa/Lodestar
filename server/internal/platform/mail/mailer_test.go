package mail_test

import (
	"bufio"
	"context"
	"encoding/base64"
	"net"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

func fakeServer(t *testing.T) (string, <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	got := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		r := bufio.NewReader(conn)
		reply := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
		reply("220 fake")
		var data strings.Builder
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			switch cmd := strings.ToUpper(strings.TrimSpace(line)); {
			case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
				reply("250 fake")
			case cmd == "DATA":
				reply("354 go")
				for {
					l, err := r.ReadString('\n')
					if err != nil || l == ".\r\n" {
						break
					}
					data.WriteString(l)
				}
				got <- data.String()
				reply("250 queued")
			case cmd == "QUIT":
				reply("221 bye")
				return
			default:
				reply("250 ok")
			}
		}
	}()
	return ln.Addr().String(), got
}

func TestSendDeliversEncodedMessage(t *testing.T) {
	addr, got := fakeServer(t)
	m, err := mail.New(mail.Config{Addr: addr, From: "Lodestar <noreply@uni.example>"})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	err = m.Send(context.Background(), mail.Message{To: "student@uni.example", Subject: "Срок сдачи", Body: "Завтра сдаём курсовую"})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	raw := <-got
	if !strings.Contains(raw, "Subject: =?utf-8?q?") || !strings.Contains(raw, "To: <student@uni.example>") {
		t.Fatalf("headers = %s", raw)
	}
	body := raw[strings.Index(raw, "\r\n\r\n")+4:]
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(body, "\r\n", ""))
	if err != nil || string(decoded) != "Завтра сдаём курсовую" {
		t.Fatalf("body = %q, %v", decoded, err)
	}
}

func TestRejectsBadConfig(t *testing.T) {
	if _, err := mail.New(mail.Config{Addr: "smtp.example", From: "a@b"}); err == nil {
		t.Fatal("address without port accepted")
	}
	if _, err := mail.New(mail.Config{Addr: "smtp.example:587", From: "not an address"}); err == nil {
		t.Fatal("bad sender accepted")
	}
}
