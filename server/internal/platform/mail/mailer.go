package mail

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

const sendTimeout = 30 * time.Second

type Mailer struct {
	addr     string
	host     string
	from     *mail.Address
	auth     smtp.Auth
	implicit bool
}

func New(cfg Config) (*Mailer, error) {
	host, port, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("mail server address: %w", err)
	}
	from, err := mail.ParseAddress(cfg.From)
	if err != nil {
		return nil, fmt.Errorf("mail sender: %w", err)
	}
	m := &Mailer{addr: cfg.Addr, host: host, from: from, implicit: port == "465"}
	if cfg.Username != "" {
		m.auth = smtp.PlainAuth("", cfg.Username, cfg.Password, host)
	}
	return m, nil
}

func (m *Mailer) Send(ctx context.Context, msg Message) error {
	to, err := mail.ParseAddress(msg.To)
	if err != nil {
		return fmt.Errorf("mail recipient: %w", err)
	}
	conn, err := m.dial(ctx)
	if err != nil {
		return err
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(sendTimeout)
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return err
	}
	c, err := smtp.NewClient(conn, m.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("mail handshake: %w", err)
	}
	defer func() { _ = c.Close() }()
	if err := m.secure(c); err != nil {
		return err
	}
	if err := c.Mail(m.from.Address); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	if err := c.Rcpt(to.Address); err != nil {
		return fmt.Errorf("mail to: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("mail data: %w", err)
	}
	msg.To = to.String()
	if _, err := w.Write(msg.encode(m.from, messageID(m.from.Address), time.Now())); err != nil {
		return fmt.Errorf("mail write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail send: %w", err)
	}
	return c.Quit()
}

func (m *Mailer) dial(ctx context.Context) (net.Conn, error) {
	d := &net.Dialer{Timeout: 10 * time.Second}
	if m.implicit {
		td := &tls.Dialer{NetDialer: d, Config: &tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}}
		return td.DialContext(ctx, "tcp", m.addr)
	}
	return d.DialContext(ctx, "tcp", m.addr)
}

func (m *Mailer) secure(c *smtp.Client) error {
	if !m.implicit {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}); err != nil {
				return fmt.Errorf("mail starttls: %w", err)
			}
		}
	}
	if m.auth != nil {
		if err := c.Auth(m.auth); err != nil {
			return fmt.Errorf("mail auth: %w", err)
		}
	}
	return nil
}

func messageID(from string) string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	domain := from[strings.LastIndexByte(from, '@')+1:]
	return hex.EncodeToString(b) + "@" + domain
}
